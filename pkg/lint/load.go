package lint

import (
	"context"
	"fmt"
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/fsutils"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type ContextLoader struct {
	cfg         *config.Config
	log         logutils.Log
	debugf      logutils.DebugFunc
	goenv       *goutil.Env
	pkgTestIDRe *regexp.Regexp
	lineCache   *fsutils.LineCache
	fileCache   *fsutils.FileCache
}

func NewContextLoader(cfg *config.Config, log logutils.Log, goenv *goutil.Env,
	lineCache *fsutils.LineCache, fileCache *fsutils.FileCache) *ContextLoader {
	return &ContextLoader{
		cfg:         cfg,
		log:         log,
		debugf:      logutils.Debug("loader"),
		goenv:       goenv,
		pkgTestIDRe: regexp.MustCompile(`^(.*) \[(.*)\.test\]`),
		lineCache:   lineCache,
		fileCache:   fileCache,
	}
}

func (cl ContextLoader) prepareBuildContext() {
	// Set GOROOT to have working cross-compilation: cross-compiled binaries
	// have invalid GOROOT. XXX: can't use runtime.GOROOT().
	goroot := cl.goenv.Get(goutil.EnvGoRoot)
	if goroot == "" {
		return
	}

	os.Setenv("GOROOT", goroot)
	build.Default.GOROOT = goroot
	build.Default.BuildTags = cl.cfg.Run.BuildTags
}

func (cl ContextLoader) makeFakeLoaderPackageInfo(pkg *packages.Package) *loader.PackageInfo {
	var errs []error
	for _, err := range pkg.Errors {
		errs = append(errs, err)
	}

	typeInfo := &types.Info{}
	if pkg.TypesInfo != nil {
		typeInfo = pkg.TypesInfo
	}

	return &loader.PackageInfo{
		Pkg:                   pkg.Types,
		Importable:            true, // not used
		TransitivelyErrorFree: !pkg.IllTyped,

		// use compiled (preprocessed) go files AST;
		// AST linters use not preprocessed go files AST
		Files:  pkg.Syntax,
		Errors: errs,
		Info:   *typeInfo,
	}
}

func (cl ContextLoader) makeFakeLoaderProgram(pkgs []*packages.Package) *loader.Program {
	var createdPkgs []*loader.PackageInfo
	for _, pkg := range pkgs {
		if pkg.IllTyped {
			// some linters crash on packages with errors,
			// skip them and warn about them in another place
			continue
		}

		pkgInfo := cl.makeFakeLoaderPackageInfo(pkg)
		createdPkgs = append(createdPkgs, pkgInfo)
	}

	allPkgs := map[*types.Package]*loader.PackageInfo{}
	for _, pkg := range createdPkgs {
		pkg := pkg
		allPkgs[pkg.Pkg] = pkg
	}
	for _, pkg := range pkgs {
		if pkg.IllTyped {
			// some linters crash on packages with errors,
			// skip them and warn about them in another place
			continue
		}

		for _, impPkg := range pkg.Imports {
			// don't use astcache for imported packages: we don't find issues in cgo imported deps
			pkgInfo := cl.makeFakeLoaderPackageInfo(impPkg)
			allPkgs[pkgInfo.Pkg] = pkgInfo
		}
	}

	return &loader.Program{
		Fset:        pkgs[0].Fset,
		Imported:    nil,         // not used without .Created in any linter
		Created:     createdPkgs, // all initial packages
		AllPackages: allPkgs,     // all initial packages and their depndencies
	}
}

func (cl ContextLoader) buildSSAProgram(pkgs []*packages.Package) *ssa.Program {
	startedAt := time.Now()
	var pkgsBuiltDuration time.Duration
	defer func() {
		cl.log.Infof("SSA repr building timing: packages building %s, total %s",
			pkgsBuiltDuration, time.Since(startedAt))
	}()

	ssaProg, _ := ssautil.Packages(pkgs, ssa.GlobalDebug)
	pkgsBuiltDuration = time.Since(startedAt)
	ssaProg.Build()
	return ssaProg
}

func (cl ContextLoader) findLoadMode(linters []*linter.Config) packages.LoadMode {
	//TODO: specify them in linters: need more fine-grained control.
	// e.g. NeedTypesSizes is needed only for go vet
	loadMode := packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	for _, lc := range linters {
		if lc.NeedsTypeInfo {
			loadMode |= packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedTypesInfo | packages.NeedSyntax
		}
		if lc.NeedsDepsTypeInfo {
			loadMode |= packages.NeedDeps
		}
	}

	return loadMode
}

func (cl ContextLoader) buildArgs() []string {
	args := cl.cfg.Run.Args
	if len(args) == 0 {
		return []string{"./..."}
	}

	var retArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, ".") || filepath.IsAbs(arg) {
			retArgs = append(retArgs, arg)
		} else {
			// go/packages doesn't work well if we don't have prefix ./ for local packages
			retArgs = append(retArgs, fmt.Sprintf(".%c%s", filepath.Separator, arg))
		}
	}

	return retArgs
}

func (cl ContextLoader) makeBuildFlags() ([]string, error) {
	var buildFlags []string

	if len(cl.cfg.Run.BuildTags) != 0 {
		// go help build
		buildFlags = append(buildFlags, "-tags", strings.Join(cl.cfg.Run.BuildTags, " "))
	}

	mod := cl.cfg.Run.ModulesDownloadMode
	if mod != "" {
		// go help modules
		allowedMods := []string{"release", "readonly", "vendor"}
		var ok bool
		for _, am := range allowedMods {
			if am == mod {
				ok = true
				break
			}
		}
		if !ok {
			return nil, fmt.Errorf("invalid modules download path %s, only (%s) allowed", mod, strings.Join(allowedMods, "|"))
		}

		buildFlags = append(buildFlags, fmt.Sprintf("-mod=%s", cl.cfg.Run.ModulesDownloadMode))
	}

	return buildFlags, nil
}

func stringifyLoadMode(mode packages.LoadMode) string {
	m := map[packages.LoadMode]string{
		packages.NeedCompiledGoFiles: "compiled_files",
		packages.NeedDeps:            "deps",
		packages.NeedExportsFile:     "exports_file",
		packages.NeedFiles:           "files",
		packages.NeedImports:         "imports",
		packages.NeedName:            "name",
		packages.NeedSyntax:          "syntax",
		packages.NeedTypes:           "types",
		packages.NeedTypesInfo:       "types_info",
		packages.NeedTypesSizes:      "types_sizes",
	}

	var flags []string
	for flag, flagStr := range m {
		if mode&flag != 0 {
			flags = append(flags, flagStr)
		}
	}

	return fmt.Sprintf("%d (%s)", mode, strings.Join(flags, "|"))
}

func (cl ContextLoader) debugPrintLoadedPackages(pkgs []*packages.Package) {
	cl.debugf("loaded %d pkgs", len(pkgs))
	for i, pkg := range pkgs {
		var syntaxFiles []string
		for _, sf := range pkg.Syntax {
			syntaxFiles = append(syntaxFiles, pkg.Fset.Position(sf.Pos()).Filename)
		}
		cl.debugf("Loaded pkg #%d: ID=%s GoFiles=%s CompiledGoFiles=%s Syntax=%s",
			i, pkg.ID, pkg.GoFiles, pkg.CompiledGoFiles, syntaxFiles)
	}
}

func (cl ContextLoader) parseLoadedPackagesErrors(pkgs []*packages.Package) error {
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			if strings.Contains(err.Msg, "no Go files") {
				return errors.Wrapf(exitcodes.ErrNoGoFiles, "package %s", pkg.PkgPath)
			}
			if strings.Contains(err.Msg, "cannot find package") {
				// when analyzing not existing directory
				return errors.Wrap(exitcodes.ErrFailure, err.Msg)
			}
		}
	}

	return nil
}

func (cl ContextLoader) loadPackages(ctx context.Context, loadMode packages.LoadMode) ([]*packages.Package, error) {
	defer func(startedAt time.Time) {
		cl.log.Infof("Go packages loading at mode %s took %s", stringifyLoadMode(loadMode), time.Since(startedAt))
	}(time.Now())

	cl.prepareBuildContext()

	buildFlags, err := cl.makeBuildFlags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make build flags for go list")
	}

	conf := &packages.Config{
		Mode:       loadMode,
		Tests:      cl.cfg.Run.AnalyzeTests,
		Context:    ctx,
		BuildFlags: buildFlags,
		Logf:       cl.debugf,
		//TODO: use fset, parsefile, overlay
	}

	args := cl.buildArgs()
	cl.debugf("Built loader args are %s", args)
	pkgs, err := packages.Load(conf, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load program with go/packages")
	}
	cl.debugPrintLoadedPackages(pkgs)

	if err := cl.parseLoadedPackagesErrors(pkgs); err != nil {
		return nil, err
	}

	return cl.filterTestMainPackages(pkgs), nil
}

func (cl ContextLoader) tryParseTestPackage(pkg *packages.Package) (name, testName string, isTest bool) {
	matches := cl.pkgTestIDRe.FindStringSubmatch(pkg.ID)
	if matches == nil {
		return "", "", false
	}

	return matches[1], matches[2], true
}

func (cl ContextLoader) filterTestMainPackages(pkgs []*packages.Package) []*packages.Package {
	var retPkgs []*packages.Package
	for _, pkg := range pkgs {
		if pkg.Name == "main" && strings.HasSuffix(pkg.PkgPath, ".test") {
			// it's an implicit testmain package
			cl.debugf("skip pkg ID=%s", pkg.ID)
			continue
		}

		retPkgs = append(retPkgs, pkg)
	}

	return retPkgs
}

func (cl ContextLoader) filterDuplicatePackages(pkgs []*packages.Package) []*packages.Package {
	packagesWithTests := map[string]bool{}
	for _, pkg := range pkgs {
		name, _, isTest := cl.tryParseTestPackage(pkg)
		if !isTest {
			continue
		}
		packagesWithTests[name] = true
	}

	cl.debugf("package with tests: %#v", packagesWithTests)

	var retPkgs []*packages.Package
	for _, pkg := range pkgs {
		_, _, isTest := cl.tryParseTestPackage(pkg)
		if !isTest && packagesWithTests[pkg.PkgPath] {
			// If tests loading is enabled,
			// for package with files a.go and a_test.go go/packages loads two packages:
			// 1. ID=".../a" GoFiles=[a.go]
			// 2. ID=".../a [.../a.test]" GoFiles=[a.go a_test.go]
			// We need only the second package, otherwise we can get warnings about unused variables/fields/functions
			// in a.go if they are used only in a_test.go.
			cl.debugf("skip pkg ID=%s because we load it with test package", pkg.ID)
			continue
		}

		retPkgs = append(retPkgs, pkg)
	}

	return retPkgs
}

func needSSA(linters []*linter.Config) bool {
	for _, lc := range linters {
		if lc.NeedsSSARepr {
			return true
		}
	}
	return false
}

//nolint:gocyclo
func (cl ContextLoader) Load(ctx context.Context, linters []*linter.Config) (*linter.Context, error) {
	loadMode := cl.findLoadMode(linters)
	pkgs, err := cl.loadPackages(ctx, loadMode)
	if err != nil {
		return nil, err
	}

	deduplicatedPkgs := cl.filterDuplicatePackages(pkgs)

	if len(deduplicatedPkgs) == 0 {
		return nil, exitcodes.ErrNoGoFiles
	}

	var prog *loader.Program
	if loadMode&packages.NeedTypes != 0 {
		prog = cl.makeFakeLoaderProgram(deduplicatedPkgs)
	}

	var ssaProg *ssa.Program
	if needSSA(linters) {
		ssaProg = cl.buildSSAProgram(deduplicatedPkgs)
	}

	astLog := cl.log.Child("astcache")
	astCache, err := astcache.LoadFromPackages(deduplicatedPkgs, astLog)
	if err != nil {
		return nil, err
	}

	ret := &linter.Context{
		Packages: deduplicatedPkgs,

		// At least `unused` linters works properly only on original (not deduplicated) packages,
		// see https://github.com/golangci/golangci-lint/pull/585.
		OriginalPackages: pkgs,

		Program:    prog,
		SSAProgram: ssaProg,
		LoaderConfig: &loader.Config{
			Cwd:   "",  // used by depguard and fallbacked to os.Getcwd
			Build: nil, // used by depguard and megacheck and fallbacked to build.Default
		},
		Cfg:       cl.cfg,
		ASTCache:  astCache,
		Log:       cl.log,
		FileCache: cl.fileCache,
		LineCache: cl.lineCache,
	}

	separateNotCompilingPackages(ret)
	return ret, nil
}

// separateNotCompilingPackages moves not compiling packages into separate slice:
// a lot of linters crash on such packages
func separateNotCompilingPackages(lintCtx *linter.Context) {
	// Separate deduplicated packages
	goodPkgs := make([]*packages.Package, 0, len(lintCtx.Packages))
	for _, pkg := range lintCtx.Packages {
		if pkg.IllTyped {
			lintCtx.NotCompilingPackages = append(lintCtx.NotCompilingPackages, pkg)
		} else {
			goodPkgs = append(goodPkgs, pkg)
		}
	}

	lintCtx.Packages = goodPkgs
	if len(lintCtx.NotCompilingPackages) != 0 {
		lintCtx.Log.Infof("Packages that do not compile: %+v", lintCtx.NotCompilingPackages)
	}

	// Separate original (not deduplicated) packages
	goodOriginalPkgs := make([]*packages.Package, 0, len(lintCtx.OriginalPackages))
	for _, pkg := range lintCtx.OriginalPackages {
		if !pkg.IllTyped {
			goodOriginalPkgs = append(goodOriginalPkgs, pkg)
		}
	}
	lintCtx.OriginalPackages = goodOriginalPkgs
}
