// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// Package lint provides the foundation for tools like gosimple.
package lint // import "github.com/golangci/go-tools/lint"

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/constant"
	"go/printer"
	"go/token"
	"go/types"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/ssa/ssautil"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

type Job struct {
	Program *Program

	checker  string
	check    string
	problems []Problem
}

type Ignore interface {
	Match(p Problem) bool
}

type LineIgnore struct {
	File    string
	Line    int
	Checks  []string
	matched bool
	pos     token.Pos
}

func (li *LineIgnore) Match(p Problem) bool {
	if p.Position.Filename != li.File || p.Position.Line != li.Line {
		return false
	}
	for _, c := range li.Checks {
		if m, _ := filepath.Match(c, p.Check); m {
			li.matched = true
			return true
		}
	}
	return false
}

func (li *LineIgnore) String() string {
	matched := "not matched"
	if li.matched {
		matched = "matched"
	}
	return fmt.Sprintf("%s:%d %s (%s)", li.File, li.Line, strings.Join(li.Checks, ", "), matched)
}

type FileIgnore struct {
	File   string
	Checks []string
}

func (fi *FileIgnore) Match(p Problem) bool {
	if p.Position.Filename != fi.File {
		return false
	}
	for _, c := range fi.Checks {
		if m, _ := filepath.Match(c, p.Check); m {
			return true
		}
	}
	return false
}

type GlobIgnore struct {
	Pattern string
	Checks  []string
}

func (gi *GlobIgnore) Match(p Problem) bool {
	if gi.Pattern != "*" {
		pkgpath := p.Package.Path()
		if strings.HasSuffix(pkgpath, "_test") {
			pkgpath = pkgpath[:len(pkgpath)-len("_test")]
		}
		name := filepath.Join(pkgpath, filepath.Base(p.Position.Filename))
		if m, _ := filepath.Match(gi.Pattern, name); !m {
			return false
		}
	}
	for _, c := range gi.Checks {
		if m, _ := filepath.Match(c, p.Check); m {
			return true
		}
	}
	return false
}

type Program struct {
	SSA  *ssa.Program
	Prog *loader.Program
	// TODO(dh): Rename to InitialPackages?
	Packages         []*Pkg
	InitialFunctions []*ssa.Function
	AllFunctions     []*ssa.Function
	Files            []*ast.File
	Info             *types.Info
	GoVersion        int

	tokenFileMap map[*token.File]*ast.File
	astFileMap   map[*ast.File]*Pkg
}

type Func func(*Job)

// Problem represents a problem in some source code.
type Problem struct {
	pos      token.Pos
	Position token.Position // position in source file
	Text     string         // the prose that describes the problem
	Check    string
	Checker  string
	Package  *types.Package
	Ignored  bool
}

func (p *Problem) String() string {
	if p.Check == "" {
		return p.Text
	}
	return fmt.Sprintf("%s (%s)", p.Text, p.Check)
}

type Checker interface {
	Name() string
	Prefix() string
	Init(*Program)
	Funcs() map[string]Func
}

// A Linter lints Go source code.
type Linter struct {
	Checker       Checker
	Ignores       []Ignore
	GoVersion     int
	ReturnIgnored bool

	automaticIgnores []Ignore
}

func (l *Linter) ignore(p Problem) bool {
	ignored := false
	for _, ig := range l.automaticIgnores {
		// We cannot short-circuit these, as we want to record, for
		// each ignore, whether it matched or not.
		if ig.Match(p) {
			ignored = true
		}
	}
	if ignored {
		// no need to execute other ignores if we've already had a
		// match.
		return true
	}
	for _, ig := range l.Ignores {
		// We can short-circuit here, as we aren't tracking any
		// information.
		if ig.Match(p) {
			return true
		}
	}

	return false
}

func (prog *Program) File(node Positioner) *ast.File {
	return prog.tokenFileMap[prog.SSA.Fset.File(node.Pos())]
}

func (j *Job) File(node Positioner) *ast.File {
	return j.Program.File(node)
}

// TODO(dh): switch to sort.Slice when Go 1.9 lands.
type byPosition struct {
	fset *token.FileSet
	ps   []Problem
}

func (ps byPosition) Len() int {
	return len(ps.ps)
}

func (ps byPosition) Less(i int, j int) bool {
	pi, pj := ps.ps[i].Position, ps.ps[j].Position

	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	}
	if pi.Line != pj.Line {
		return pi.Line < pj.Line
	}
	if pi.Column != pj.Column {
		return pi.Column < pj.Column
	}

	return ps.ps[i].Text < ps.ps[j].Text
}

func (ps byPosition) Swap(i int, j int) {
	ps.ps[i], ps.ps[j] = ps.ps[j], ps.ps[i]
}

func parseDirective(s string) (cmd string, args []string) {
	if !strings.HasPrefix(s, "//lint:") {
		return "", nil
	}
	s = strings.TrimPrefix(s, "//lint:")
	fields := strings.Split(s, " ")
	return fields[0], fields[1:]
}

func (l *Linter) Lint(lprog *loader.Program, conf *loader.Config, ssaprog *ssa.Program) []Problem {
	pkgMap := map[*ssa.Package]*Pkg{}
	var pkgs []*Pkg
	for _, pkginfo := range lprog.InitialPackages() {
		ssapkg := ssaprog.Package(pkginfo.Pkg)
		var bp *build.Package
		if len(pkginfo.Files) != 0 {
			path := lprog.Fset.Position(pkginfo.Files[0].Pos()).Filename
			dir := filepath.Dir(path)
			var err error
			ctx := conf.Build
			if ctx == nil {
				ctx = &build.Default
			}
			bp, err = ctx.ImportDir(dir, 0)
			if err != nil {
				// shouldn't happen
			}
		}
		pkg := &Pkg{
			Package:  ssapkg,
			Info:     pkginfo,
			BuildPkg: bp,
		}
		pkgMap[ssapkg] = pkg
		pkgs = append(pkgs, pkg)
	}
	prog := &Program{
		SSA:          ssaprog,
		Prog:         lprog,
		Packages:     pkgs,
		Info:         &types.Info{},
		GoVersion:    l.GoVersion,
		tokenFileMap: map[*token.File]*ast.File{},
		astFileMap:   map[*ast.File]*Pkg{},
	}

	initial := map[*types.Package]struct{}{}
	for _, pkg := range pkgs {
		initial[pkg.Info.Pkg] = struct{}{}
	}
	for fn := range ssautil.AllFunctions(ssaprog) {
		if fn.Pkg == nil {
			continue
		}
		prog.AllFunctions = append(prog.AllFunctions, fn)
		if _, ok := initial[fn.Pkg.Pkg]; ok {
			prog.InitialFunctions = append(prog.InitialFunctions, fn)
		}
	}
	for _, pkg := range pkgs {
		prog.Files = append(prog.Files, pkg.Info.Files...)

		ssapkg := ssaprog.Package(pkg.Info.Pkg)
		for _, f := range pkg.Info.Files {
			prog.astFileMap[f] = pkgMap[ssapkg]
		}
	}

	for _, pkginfo := range lprog.AllPackages {
		for _, f := range pkginfo.Files {
			tf := lprog.Fset.File(f.Pos())
			prog.tokenFileMap[tf] = f
		}
	}

	var out []Problem
	l.automaticIgnores = nil
	for _, pkginfo := range lprog.InitialPackages() {
		for _, f := range pkginfo.Files {
			cm := ast.NewCommentMap(lprog.Fset, f, f.Comments)
			for node, cgs := range cm {
				for _, cg := range cgs {
					for _, c := range cg.List {
						if !strings.HasPrefix(c.Text, "//lint:") {
							continue
						}
						cmd, args := parseDirective(c.Text)
						switch cmd {
						case "ignore", "file-ignore":
							if len(args) < 2 {
								// FIXME(dh): this causes duplicated warnings when using megacheck
								p := Problem{
									pos:      c.Pos(),
									Position: prog.DisplayPosition(c.Pos()),
									Text:     "malformed linter directive; missing the required reason field?",
									Check:    "",
									Checker:  l.Checker.Name(),
									Package:  nil,
								}
								out = append(out, p)
								continue
							}
						default:
							// unknown directive, ignore
							continue
						}
						checks := strings.Split(args[0], ",")
						pos := prog.DisplayPosition(node.Pos())
						var ig Ignore
						switch cmd {
						case "ignore":
							ig = &LineIgnore{
								File:   pos.Filename,
								Line:   pos.Line,
								Checks: checks,
								pos:    c.Pos(),
							}
						case "file-ignore":
							ig = &FileIgnore{
								File:   pos.Filename,
								Checks: checks,
							}
						}
						l.automaticIgnores = append(l.automaticIgnores, ig)
					}
				}
			}
		}
	}

	sizes := struct {
		types      int
		defs       int
		uses       int
		implicits  int
		selections int
		scopes     int
	}{}
	for _, pkg := range pkgs {
		sizes.types += len(pkg.Info.Info.Types)
		sizes.defs += len(pkg.Info.Info.Defs)
		sizes.uses += len(pkg.Info.Info.Uses)
		sizes.implicits += len(pkg.Info.Info.Implicits)
		sizes.selections += len(pkg.Info.Info.Selections)
		sizes.scopes += len(pkg.Info.Info.Scopes)
	}
	prog.Info.Types = make(map[ast.Expr]types.TypeAndValue, sizes.types)
	prog.Info.Defs = make(map[*ast.Ident]types.Object, sizes.defs)
	prog.Info.Uses = make(map[*ast.Ident]types.Object, sizes.uses)
	prog.Info.Implicits = make(map[ast.Node]types.Object, sizes.implicits)
	prog.Info.Selections = make(map[*ast.SelectorExpr]*types.Selection, sizes.selections)
	prog.Info.Scopes = make(map[ast.Node]*types.Scope, sizes.scopes)
	for _, pkg := range pkgs {
		for k, v := range pkg.Info.Info.Types {
			prog.Info.Types[k] = v
		}
		for k, v := range pkg.Info.Info.Defs {
			prog.Info.Defs[k] = v
		}
		for k, v := range pkg.Info.Info.Uses {
			prog.Info.Uses[k] = v
		}
		for k, v := range pkg.Info.Info.Implicits {
			prog.Info.Implicits[k] = v
		}
		for k, v := range pkg.Info.Info.Selections {
			prog.Info.Selections[k] = v
		}
		for k, v := range pkg.Info.Info.Scopes {
			prog.Info.Scopes[k] = v
		}
	}
	l.Checker.Init(prog)

	funcs := l.Checker.Funcs()
	var keys []string
	for k := range funcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var jobs []*Job
	for _, k := range keys {
		j := &Job{
			Program: prog,
			checker: l.Checker.Name(),
			check:   k,
		}
		jobs = append(jobs, j)
	}
	wg := &sync.WaitGroup{}
	crashesMap := sync.Map{}
	for i, j := range jobs {
		wg.Add(1)
		go func(i int, j *Job) {
			defer func() {
				if err := recover(); err != nil {
					crashesMap.Store(i, err)
				}
			}()
			defer wg.Done()
			fn := funcs[j.check]
			if fn == nil {
				return
			}
			fn(j)
		}(i, j)
	}
	wg.Wait()

	for i := range jobs {
		err, ok := crashesMap.Load(i)
		if !ok {
			continue
		}
		panic(err) // restore panic but to main goroutine to be properly catched
	}

	for _, j := range jobs {
		for _, p := range j.problems {
			p.Ignored = l.ignore(p)
			if l.ReturnIgnored || !p.Ignored {
				out = append(out, p)
			}
		}
	}

	for _, ig := range l.automaticIgnores {
		ig, ok := ig.(*LineIgnore)
		if !ok {
			continue
		}
		if ig.matched {
			continue
		}
		for _, c := range ig.Checks {
			idx := strings.IndexFunc(c, func(r rune) bool {
				return unicode.IsNumber(r)
			})
			if idx == -1 {
				// malformed check name, backing out
				continue
			}
			if c[:idx] != l.Checker.Prefix() {
				// not for this checker
				continue
			}
			p := Problem{
				pos:      ig.pos,
				Position: prog.DisplayPosition(ig.pos),
				Text:     "this linter directive didn't match anything; should it be removed?",
				Check:    "",
				Checker:  l.Checker.Name(),
				Package:  nil,
			}
			out = append(out, p)
		}
	}

	sort.Sort(byPosition{lprog.Fset, out})
	return out
}

// Pkg represents a package being linted.
type Pkg struct {
	*ssa.Package
	Info     *loader.PackageInfo
	BuildPkg *build.Package
}

func (p Pkg) GetPackage() *types.Package {
	if p.Package == nil {
		return nil
	}

	return p.Pkg
}

type packager interface {
	Package() *ssa.Package
}

func IsExample(fn *ssa.Function) bool {
	if !strings.HasPrefix(fn.Name(), "Example") {
		return false
	}
	f := fn.Prog.Fset.File(fn.Pos())
	if f == nil {
		return false
	}
	return strings.HasSuffix(f.Name(), "_test.go")
}

func (j *Job) IsInTest(node Positioner) bool {
	f := j.Program.SSA.Fset.File(node.Pos())
	return f != nil && strings.HasSuffix(f.Name(), "_test.go")
}

func (j *Job) IsInMain(node Positioner) bool {
	if node, ok := node.(packager); ok {
		return node.Package().Pkg.Name() == "main"
	}
	pkg := j.NodePackage(node)
	if pkg == nil {
		return false
	}
	return pkg.Pkg.Name() == "main"
}

type Positioner interface {
	Pos() token.Pos
}

func (prog *Program) DisplayPosition(p token.Pos) token.Position {
	// The //line compiler directive can be used to change the file
	// name and line numbers associated with code. This can, for
	// example, be used by code generation tools. The most prominent
	// example is 'go tool cgo', which uses //line directives to refer
	// back to the original source code.
	//
	// In the context of our linters, we need to treat these
	// directives differently depending on context. For cgo files, we
	// want to honour the directives, so that line numbers are
	// adjusted correctly. For all other files, we want to ignore the
	// directives, so that problems are reported at their actual
	// position and not, for example, a yacc grammar file. This also
	// affects the ignore mechanism, since it operates on the position
	// information stored within problems. With this implementation, a
	// user will ignore foo.go, not foo.y

	pkg := prog.astFileMap[prog.tokenFileMap[prog.Prog.Fset.File(p)]]
	bp := pkg.BuildPkg
	adjPos := prog.Prog.Fset.Position(p)
	if bp == nil {
		// couldn't find the package for some reason (deleted? faulty
		// file system?)
		return adjPos
	}
	base := filepath.Base(adjPos.Filename)
	for _, f := range bp.CgoFiles {
		if f == base {
			// this is a cgo file, use the adjusted position
			return adjPos
		}
	}
	// not a cgo file, ignore //line directives
	return prog.Prog.Fset.PositionFor(p, false)
}

func (j *Job) Errorf(n Positioner, format string, args ...interface{}) *Problem {
	tf := j.Program.SSA.Fset.File(n.Pos())
	f := j.Program.tokenFileMap[tf]
	pkg := j.Program.astFileMap[f].GetPackage()

	pos := j.Program.DisplayPosition(n.Pos())
	problem := Problem{
		pos:      n.Pos(),
		Position: pos,
		Text:     fmt.Sprintf(format, args...),
		Check:    j.check,
		Checker:  j.checker,
		Package:  pkg,
	}
	j.problems = append(j.problems, problem)
	return &j.problems[len(j.problems)-1]
}

func (j *Job) Render(x interface{}) string {
	fset := j.Program.SSA.Fset
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func (j *Job) RenderArgs(args []ast.Expr) string {
	var ss []string
	for _, arg := range args {
		ss = append(ss, j.Render(arg))
	}
	return strings.Join(ss, ", ")
}

func IsIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// isBlank returns whether id is the blank identifier "_".
// If id == nil, the answer is false.
func IsBlank(id ast.Expr) bool {
	ident, ok := id.(*ast.Ident)
	return ok && ident.Name == "_"
}

func IsZero(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == "0"
}

func (j *Job) IsNil(expr ast.Expr) bool {
	return j.Program.Info.Types[expr].IsNil()
}

func (j *Job) BoolConst(expr ast.Expr) bool {
	val := j.Program.Info.ObjectOf(expr.(*ast.Ident)).(*types.Const).Val()
	return constant.BoolVal(val)
}

func (j *Job) IsBoolConst(expr ast.Expr) bool {
	// We explicitly don't support typed bools because more often than
	// not, custom bool types are used as binary enums and the
	// explicit comparison is desired.

	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	obj := j.Program.Info.ObjectOf(ident)
	c, ok := obj.(*types.Const)
	if !ok {
		return false
	}
	basic, ok := c.Type().(*types.Basic)
	if !ok {
		return false
	}
	if basic.Kind() != types.UntypedBool && basic.Kind() != types.Bool {
		return false
	}
	return true
}

func (j *Job) ExprToInt(expr ast.Expr) (int64, bool) {
	tv := j.Program.Info.Types[expr]
	if tv.Value == nil {
		return 0, false
	}
	if tv.Value.Kind() != constant.Int {
		return 0, false
	}
	return constant.Int64Val(tv.Value)
}

func (j *Job) ExprToString(expr ast.Expr) (string, bool) {
	val := j.Program.Info.Types[expr].Value
	if val == nil {
		return "", false
	}
	if val.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(val), true
}

func (j *Job) NodePackage(node Positioner) *Pkg {
	f := j.File(node)
	return j.Program.astFileMap[f]
}

func IsGenerated(f *ast.File) bool {
	comments := f.Comments
	if len(comments) > 0 {
		comment := comments[0].Text()
		return strings.Contains(comment, "Code generated by") ||
			strings.Contains(comment, "DO NOT EDIT")
	}
	return false
}

func Preamble(f *ast.File) string {
	cutoff := f.Package
	if f.Doc != nil {
		cutoff = f.Doc.Pos()
	}
	var out []string
	for _, cmt := range f.Comments {
		if cmt.Pos() >= cutoff {
			break
		}
		out = append(out, cmt.Text())
	}
	return strings.Join(out, "\n")
}

func IsPointerLike(T types.Type) bool {
	switch T := T.Underlying().(type) {
	case *types.Interface, *types.Chan, *types.Map, *types.Pointer:
		return true
	case *types.Basic:
		return T.Kind() == types.UnsafePointer
	}
	return false
}

func (j *Job) IsGoVersion(minor int) bool {
	return j.Program.GoVersion >= minor
}

func (j *Job) IsCallToAST(node ast.Node, name string) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	fn, ok := j.Program.Info.ObjectOf(sel.Sel).(*types.Func)
	return ok && fn.FullName() == name
}

func (j *Job) IsCallToAnyAST(node ast.Node, names ...string) bool {
	for _, name := range names {
		if j.IsCallToAST(node, name) {
			return true
		}
	}
	return false
}

func (j *Job) SelectorName(expr *ast.SelectorExpr) string {
	sel := j.Program.Info.Selections[expr]
	if sel == nil {
		if x, ok := expr.X.(*ast.Ident); ok {
			pkg, ok := j.Program.Info.ObjectOf(x).(*types.PkgName)
			if !ok {
				// This shouldn't happen
				return fmt.Sprintf("%s.%s", x.Name, expr.Sel.Name)
			}
			return fmt.Sprintf("%s.%s", pkg.Imported().Path(), expr.Sel.Name)
		}
		panic(fmt.Sprintf("unsupported selector: %v", expr))
	}
	return fmt.Sprintf("(%s).%s", sel.Recv(), sel.Obj().Name())
}

func CallName(call *ssa.CallCommon) string {
	if call.IsInvoke() {
		return ""
	}
	switch v := call.Value.(type) {
	case *ssa.Function:
		fn, ok := v.Object().(*types.Func)
		if !ok {
			return ""
		}
		return fn.FullName()
	case *ssa.Builtin:
		return v.Name()
	}
	return ""
}

func IsCallTo(call *ssa.CallCommon, name string) bool {
	return CallName(call) == name
}

func FilterDebug(instr []ssa.Instruction) []ssa.Instruction {
	var out []ssa.Instruction
	for _, ins := range instr {
		if _, ok := ins.(*ssa.DebugRef); !ok {
			out = append(out, ins)
		}
	}
	return out
}

func NodeFns(pkgs []*Pkg) map[ast.Node]*ssa.Function {
	out := map[ast.Node]*ssa.Function{}

	wg := &sync.WaitGroup{}
	chNodeFns := make(chan map[ast.Node]*ssa.Function, runtime.NumCPU()*2)
	for _, pkg := range pkgs {
		if pkg.Package == nil { // package wasn't loaded, probably it doesn't compile
			continue
		}
		pkg := pkg
		wg.Add(1)
		go func() {
			m := map[ast.Node]*ssa.Function{}
			for _, f := range pkg.Info.Files {
				ast.Walk(&globalVisitor{m, pkg, f}, f)
			}
			chNodeFns <- m
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(chNodeFns)
	}()

	for nodeFns := range chNodeFns {
		for k, v := range nodeFns {
			out[k] = v
		}
	}

	return out
}

type globalVisitor struct {
	m   map[ast.Node]*ssa.Function
	pkg *Pkg
	f   *ast.File
}

func (v *globalVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.CallExpr:
		v.m[node] = v.pkg.Func("init")
		return v
	case *ast.FuncDecl, *ast.FuncLit:
		nv := &fnVisitor{v.m, v.f, v.pkg, nil}
		return nv.Visit(node)
	default:
		return v
	}
}

type fnVisitor struct {
	m     map[ast.Node]*ssa.Function
	f     *ast.File
	pkg   *Pkg
	ssafn *ssa.Function
}

func (v *fnVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.FuncDecl:
		var ssafn *ssa.Function
		ssafn = v.pkg.Prog.FuncValue(v.pkg.Info.ObjectOf(node.Name).(*types.Func))
		v.m[node] = ssafn
		if ssafn == nil {
			return nil
		}
		return &fnVisitor{v.m, v.f, v.pkg, ssafn}
	case *ast.FuncLit:
		var ssafn *ssa.Function
		path, _ := astutil.PathEnclosingInterval(v.f, node.Pos(), node.Pos())
		ssafn = ssa.EnclosingFunction(v.pkg.Package, path)
		v.m[node] = ssafn
		if ssafn == nil {
			return nil
		}
		return &fnVisitor{v.m, v.f, v.pkg, ssafn}
	case nil:
		return nil
	default:
		v.m[node] = v.ssafn
		return v
	}
}
