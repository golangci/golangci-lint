package wsl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"strings"
)

type Configuration struct {
	// StrictAppend will do strict checking when assigning from append (x =
	// append(x, y)). If this is set to true the append call must append either
	// a variable assigned, called or used on the line above. Example on not
	// allowed when this is true:
	//
	//  x := []string{}
	//  y := "not going in X"
	//  x = append(x, "not y") // This is not allowed with StrictAppend
	//  z := "going in X"
	//
	//  x = append(x, z) // This is allowed with StrictAppend
	//
	//  m := transform(z)
	//  x = append(x, z) // So is this because Z is used above.
	StrictAppend bool

	// AllowAssignAndCallCuddle allows assignments to be cuddled with variables
	// used in calls on line above and calls to be cuddled with assignments of
	// variables used in call on line above.
	// Example supported with this set to true:
	//
	//  x.Call()
	//  x = Assign()
	//  x.AnotherCall()
	//  x = AnotherAssign()
	AllowAssignAndCallCuddle bool

	// AllowMultiLineAssignCuddle allows cuddling to assignments even if they
	// span over multiple lines. This defaults to true which allows the
	// following example:
	//
	//  err := function(
	//  	"multiple", "lines",
	//  )
	//  if err != nil {
	//  	// ...
	//  }
	AllowMultiLineAssignCuddle bool

	// If the number of lines in a case block is equal to or lager than this
	// number, the case *must* end white a newline.
	CaseForceTrailingWhitespaceLimit int

	// AllowTrailingComment will allow blocks to end with comments.
	AllowTrailingComment bool

	// AllowCuddleDeclaration will allow multiple var/declaration statements to
	// be cuddled. This defaults to false but setting it to true will enable the
	// following example:
	//  var foo bool
	//  var err error
	AllowCuddleDeclaration bool

	// AllowCuddleWithCalls is a list of call idents that everything can be
	// cuddled with. Defaults to calls looking like locks to support a flow like
	// this:
	//
	//  mu.Lock()
	//  allow := thisAssignment
	AllowCuddleWithCalls []string

	// AllowCuddleWithRHS is a list of right hand side variables that is allowed
	// to be cuddled with anything. Defaults to assignments or calls looking
	// like unlocks to support a flow like this:
	//
	//  allow := thisAssignment()
	//  mu.Unlock()
	AllowCuddleWithRHS []string
}

// DefaultConfig returns default configuration
func DefaultConfig() Configuration {
	return Configuration{
		StrictAppend:                     true,
		AllowAssignAndCallCuddle:         true,
		AllowMultiLineAssignCuddle:       true,
		AllowTrailingComment:             false,
		CaseForceTrailingWhitespaceLimit: 0,
		AllowCuddleWithCalls:             []string{"Lock", "RLock"},
		AllowCuddleWithRHS:               []string{"Unlock", "RUnlock"},
	}
}

// Result represents the result of one error.
type Result struct {
	FileName   string
	LineNumber int
	Position   token.Position
	Reason     string
}

// String returns the filename, line number and reason of a Result.
func (r *Result) String() string {
	return fmt.Sprintf("%s:%d: %s", r.FileName, r.LineNumber, r.Reason)
}

type Processor struct {
	config   Configuration
	result   []Result
	warnings []string
	fileSet  *token.FileSet
	file     *ast.File
}

// NewProcessor will create a Processor.
func NewProcessorWithConfig(cfg Configuration) *Processor {
	return &Processor{
		result: []Result{},
		config: cfg,
	}
}

// NewProcessor will create a Processor.
func NewProcessor() *Processor {
	return NewProcessorWithConfig(DefaultConfig())
}

// ProcessFiles takes a string slice with file names (full paths) and lints
// them.
// nolint: gocritic
func (p *Processor) ProcessFiles(filenames []string) ([]Result, []string) {
	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		p.process(filename, data)
	}

	return p.result, p.warnings
}

func (p *Processor) process(filename string, data []byte) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, data, parser.ParseComments)

	// If the file is not parsable let's add a syntax error and move on.
	if err != nil {
		p.result = append(p.result, Result{
			FileName:   filename,
			LineNumber: 0,
			Reason:     fmt.Sprintf("invalid syntax, file cannot be linted (%s)", err.Error()),
		})

		return
	}

	p.fileSet = fileSet
	p.file = file

	for _, d := range p.file.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			p.parseBlockBody(v.Name, v.Body)
		case *ast.GenDecl:
			// `go fmt` will handle proper spacing for GenDecl such as imports,
			// constants etc.
		default:
			p.addWarning("type not implemented", d.Pos(), v)
		}
	}
}

// parseBlockBody will parse any kind of block statements such as switch cases
// and if statements. A list of Result is returned.
func (p *Processor) parseBlockBody(ident *ast.Ident, block *ast.BlockStmt) {
	// Nothing to do if there's no value.
	if reflect.ValueOf(block).IsNil() {
		return
	}

	// Start by finding leading and trailing whitespaces.
	p.findLeadingAndTrailingWhitespaces(ident, block, nil)

	// Parse the block body contents.
	p.parseBlockStatements(block.List)
}

// parseBlockStatements will parse all the statements found in the body of a
// node. A list of Result is returned.
// nolint: gocognit
func (p *Processor) parseBlockStatements(statements []ast.Stmt) {
	for i, stmt := range statements {
		// Start by checking if this statement is another block (other than if,
		// for and range). This could be assignment to a function, defer or go
		// call with an inline function or similar. If this is found we start by
		// parsing this body block before moving on.
		for _, stmtBlocks := range p.findBlockStmt(stmt) {
			p.parseBlockBody(nil, stmtBlocks)
		}

		firstBodyStatement := p.firstBodyStatement(i, statements)

		// First statement, nothing to do.
		if i == 0 {
			continue
		}

		previousStatement := statements[i-1]

		// If the last statement didn't end one line above the current statement
		// we know we're not cuddled so just move on.
		if p.nodeEnd(previousStatement) != p.nodeStart(stmt)-1 {
			continue
		}

		// We know we're cuddled, extract assigned variables on the line above
		// which is the only thing we allow cuddling with. If the assignment is
		// made over multiple lines we should not allow cuddling.
		var assignedOnLineAbove []string

		// We want to keep track of what was called on the line above to support
		// special handling of things such as mutexes.
		var calledOnLineAbove []string

		// Check if the previous statement spans over multiple lines.
		var isMultiLineAssignment = p.nodeStart(previousStatement) != p.nodeStart(stmt)-1

		// Ensure previous line is not a multi line assignment and if not get
		// rightAndLeftHandSide assigned variables.
		if !isMultiLineAssignment {
			assignedOnLineAbove = p.findLHS(previousStatement)
			calledOnLineAbove = p.findRHS(previousStatement)
		}

		// If previous assignment is multi line and we allow it, fetch
		// assignments (but only assignments).
		if isMultiLineAssignment && p.config.AllowMultiLineAssignCuddle {
			if _, ok := previousStatement.(*ast.AssignStmt); ok {
				assignedOnLineAbove = p.findLHS(previousStatement)
			}
		}

		// We could potentially have a block which require us to check the first
		// argument before ruling out an allowed cuddle.
		var assignedFirstInBlock []string

		if firstBodyStatement != nil {
			assignedFirstInBlock = p.findLHS(firstBodyStatement)
		}

		var (
			leftHandSide                = p.findLHS(stmt)
			rightHandSide               = p.findRHS(stmt)
			rightAndLeftHandSide        = append(leftHandSide, rightHandSide...)
			calledOrAssignedOnLineAbove = append(calledOnLineAbove, assignedOnLineAbove...)
		)

		// If we called some kind of lock on the line above we allow cuddling
		// anything.
		if atLeastOneInListsMatch(calledOnLineAbove, p.config.AllowCuddleWithCalls) {
			continue
		}

		// If we call some kind of unlock on this line we allow cuddling with
		// anything.
		if atLeastOneInListsMatch(rightHandSide, p.config.AllowCuddleWithRHS) {
			continue
		}

		moreThanOneStatementAbove := func() bool {
			if i < 2 {
				return false
			}

			statementBeforePreviousStatement := statements[i-2]

			return p.nodeStart(previousStatement)-1 == p.nodeEnd(statementBeforePreviousStatement)
		}

		isLastStatementInBlockOfOnlyTwoLines := func() bool {
			// If we're the last statement, check if there's no more than two
			// lines from the starting statement and the end of this statement.
			// This is to support short return functions such as:
			// func (t *Typ) X() {
			//     t.X = true
			//     return t
			// }
			// nolint: gocritic
			if i == len(statements)-1 && i == 1 {
				if p.nodeEnd(stmt)-p.nodeStart(previousStatement) <= 2 {
					return true
				}
			}

			return false
		}

		switch t := stmt.(type) {
		case *ast.IfStmt:
			if len(assignedOnLineAbove) == 0 {
				p.addError(t.Pos(), "if statements should only be cuddled with assignments")
				continue
			}

			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before if statement")
				continue
			}

			if atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				continue
			}

			if atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
				continue
			}

			p.addError(t.Pos(), "if statements should only be cuddled with assignments used in the if statement itself")
		case *ast.ReturnStmt:
			if isLastStatementInBlockOfOnlyTwoLines() {
				continue
			}

			p.addError(t.Pos(), "return statements should not be cuddled if block has more than two lines")
		case *ast.BranchStmt:
			if isLastStatementInBlockOfOnlyTwoLines() {
				continue
			}

			p.addError(t.Pos(), "branch statements should not be cuddled if block has more than two lines")
		case *ast.AssignStmt:
			// append is usually an assignment but should not be allowed to be
			// cuddled with anything not appended.
			if len(rightHandSide) > 0 && rightHandSide[len(rightHandSide)-1] == "append" {
				if p.config.StrictAppend {
					if !atLeastOneInListsMatch(calledOrAssignedOnLineAbove, rightHandSide) {
						p.addError(t.Pos(), "append only allowed to cuddle with appended value")
					}
				}

				continue
			}

			if _, ok := previousStatement.(*ast.AssignStmt); ok {
				continue
			}

			// If the assignment is from a type or variable called on the line
			// above we can allow it by setting AllowAssignAndCallCuddle to
			// true.
			// Example (x is used):
			//  x.function()
			//  a.Field = x.anotherFunction()
			if p.config.AllowAssignAndCallCuddle {
				if atLeastOneInListsMatch(calledOrAssignedOnLineAbove, rightAndLeftHandSide) {
					continue
				}
			}

			p.addError(t.Pos(), "assignments should only be cuddled with other assignments")
		case *ast.DeclStmt:
			if !p.config.AllowCuddleDeclaration {
				p.addError(t.Pos(), "declarations should never be cuddled")
			}
		case *ast.ExprStmt:
			switch previousStatement.(type) {
			case *ast.DeclStmt, *ast.ReturnStmt:
				p.addError(t.Pos(), "expressions should not be cuddled with declarations or returns")
			case *ast.IfStmt, *ast.RangeStmt, *ast.SwitchStmt:
				p.addError(t.Pos(), "expressions should not be cuddled with blocks")
			}

			// If the expression is called on a type or variable used or
			// assigned on the line we can allow it by setting
			// AllowAssignAndCallCuddle to true.
			// Example of allowed cuddled (x is used):
			//  a.Field = x.func()
			//  x.function()
			if p.config.AllowAssignAndCallCuddle {
				if atLeastOneInListsMatch(calledOrAssignedOnLineAbove, rightAndLeftHandSide) {
					continue
				}
			}

			// If we assigned variables on the line above but didn't use them in
			// this expression there should probably be a newline between them.
			if len(assignedOnLineAbove) > 0 && !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				p.addError(t.Pos(), "only cuddled expressions if assigning variable or using from line above")
			}
		case *ast.RangeStmt:
			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before range statement")
				continue
			}

			if !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					p.addError(t.Pos(), "ranges should only be cuddled with assignments used in the iteration")
				}
			}
		case *ast.DeferStmt:
			if _, ok := previousStatement.(*ast.DeferStmt); ok {
				// We may cuddle multiple defers to group logic.
				continue
			}

			// Special treatment of deferring body closes after error checking
			// according to best practices. See
			// https://github.com/bombsimon/wsl/issues/31 which links to
			// discussion about error handling after HTTP requests. This is hard
			// coded and very specific but for now this is to be seen as a
			// special case. What this does is that it *only* allows a defer
			// statement with `Close` on the right hand side to be cuddled with
			// an if-statement to support this:
			//  resp, err := client.Do(req)
			//  if err != nil {
			//      return err
			//  }
			//  defer resp.Body.Close()
			if _, ok := previousStatement.(*ast.IfStmt); ok {
				if atLeastOneInListsMatch(rightHandSide, []string{"Close"}) {
					continue
				}
			}

			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before defer statement")

				continue
			}

			// Be extra nice with RHS, it's common to use this for locks:
			// m.Lock()
			// defer m.Unlock()
			previousRHS := p.findRHS(previousStatement)
			if atLeastOneInListsMatch(rightHandSide, previousRHS) {
				continue
			}

			if !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				p.addError(t.Pos(), "defer statements should only be cuddled with expressions on same variable")
			}
		case *ast.ForStmt:
			if len(rightAndLeftHandSide) == 0 {
				p.addError(t.Pos(), "for statement without condition should never be cuddled")

				continue
			}

			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before for statement")

				continue
			}

			// The same rule applies for ranges as for if statements, see
			// comments regarding variable usages on the line before or as the
			// first line in the block for details.
			if !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					p.addError(t.Pos(), "for statements should only be cuddled with assignments used in the iteration")
				}
			}
		case *ast.GoStmt:
			if _, ok := previousStatement.(*ast.GoStmt); ok {
				continue
			}

			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before go statement")

				continue
			}

			if !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				p.addError(t.Pos(), "go statements can only invoke functions assigned on line above")
			}
		case *ast.SwitchStmt:
			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before switch statement")

				continue
			}

			if !atLeastOneInListsMatch(rightAndLeftHandSide, assignedOnLineAbove) {
				if len(rightAndLeftHandSide) == 0 {
					p.addError(t.Pos(), "anonymous switch statements should never be cuddled")
				} else {
					p.addError(t.Pos(), "switch statements should only be cuddled with variables switched")
				}
			}
		case *ast.TypeSwitchStmt:
			if moreThanOneStatementAbove() {
				p.addError(t.Pos(), "only one cuddle assignment allowed before type switch statement")

				continue
			}

			// Allowed to type assert on variable assigned on line above.
			if !atLeastOneInListsMatch(rightHandSide, assignedOnLineAbove) {
				// Allow type assertion on variables used in the first case
				// immediately.
				if !atLeastOneInListsMatch(assignedOnLineAbove, assignedFirstInBlock) {
					p.addError(t.Pos(), "type switch statements should only be cuddled with variables switched")
				}
			}
		case *ast.CaseClause, *ast.CommClause:
			// Case clauses will be checked by not allowing leading ot trailing
			// whitespaces within the block. There's nothing in the case itself
			// that may be cuddled.
		default:
			p.addWarning("stmt type not implemented", t.Pos(), t)
		}
	}
}

// firstBodyStatement returns the first statement inside a body block. This is
// because variables may be cuddled with conditions or statements if it's used
// directly as the first argument inside a body.
// The body will then be parsed as a *ast.BlockStmt (regular block) or as a list
// of []ast.Stmt (case block).
func (p *Processor) firstBodyStatement(i int, allStmt []ast.Stmt) ast.Node {
	stmt := allStmt[i]

	// Start by checking if the statement has a body (probably if-statement,
	// a range, switch case or similar. Whenever a body is found we start by
	// parsing it before moving on in the AST.
	statementBody := reflect.Indirect(reflect.ValueOf(stmt)).FieldByName("Body")

	// Some cases allow cuddling depending on the first statement in a body
	// of a block or case. If possible extract the first statement.
	var firstBodyStatement ast.Node

	if !statementBody.IsValid() {
		return firstBodyStatement
	}

	switch statementBodyContent := statementBody.Interface().(type) {
	case *ast.BlockStmt:
		if len(statementBodyContent.List) > 0 {
			firstBodyStatement = statementBodyContent.List[0]

			// If the first body statement is a *ast.CaseClause we're
			// actually interested in the **next** body to know what's
			// inside the first case.
			if x, ok := firstBodyStatement.(*ast.CaseClause); ok {
				if len(x.Body) > 0 {
					firstBodyStatement = x.Body[0]
				}
			}
		}

		p.parseBlockBody(nil, statementBodyContent)
	case []ast.Stmt:
		// The Body field for an *ast.CaseClause or *ast.CommClause is of type
		// []ast.Stmt. We must check leading and trailing whitespaces and then
		// pass the statements to parseBlockStatements to parse it's content.
		var nextStatement ast.Node

		// Check if there's more statements (potential cases) after the
		// current one.
		if len(allStmt)-1 > i {
			nextStatement = allStmt[i+1]
		}

		p.findLeadingAndTrailingWhitespaces(nil, stmt, nextStatement)
		p.parseBlockStatements(statementBodyContent)
	default:
		p.addWarning(
			"body statement type not implemented ",
			stmt.Pos(), statementBodyContent,
		)
	}

	return firstBodyStatement
}

func (p *Processor) findLHS(node ast.Node) []string {
	var lhs []string

	if node == nil {
		return lhs
	}

	switch t := node.(type) {
	case *ast.BasicLit, *ast.FuncLit, *ast.SelectStmt,
		*ast.LabeledStmt, *ast.ForStmt, *ast.SwitchStmt,
		*ast.ReturnStmt, *ast.GoStmt, *ast.CaseClause,
		*ast.CommClause, *ast.CallExpr, *ast.UnaryExpr,
		*ast.BranchStmt, *ast.TypeSpec, *ast.ChanType,
		*ast.DeferStmt, *ast.TypeAssertExpr, *ast.RangeStmt:
	// Nothing to add to LHS
	case *ast.IncDecStmt:
		return p.findLHS(t.X)
	case *ast.Ident:
		return []string{t.Name}
	case *ast.AssignStmt:
		for _, v := range t.Lhs {
			lhs = append(lhs, p.findLHS(v)...)
		}
	case *ast.GenDecl:
		for _, v := range t.Specs {
			lhs = append(lhs, p.findLHS(v)...)
		}
	case *ast.ValueSpec:
		for _, v := range t.Names {
			lhs = append(lhs, p.findLHS(v)...)
		}
	case *ast.BlockStmt:
		for _, v := range t.List {
			lhs = append(lhs, p.findLHS(v)...)
		}
	case *ast.BinaryExpr:
		return append(
			p.findLHS(t.X),
			p.findLHS(t.Y)...,
		)
	case *ast.DeclStmt:
		return p.findLHS(t.Decl)
	case *ast.IfStmt:
		return p.findLHS(t.Cond)
	case *ast.TypeSwitchStmt:
		return p.findLHS(t.Assign)
	case *ast.SendStmt:
		return p.findLHS(t.Chan)
	default:
		if x, ok := maybeX(t); ok {
			return p.findLHS(x)
		}

		p.addWarning("UNKNOWN LHS", t.Pos(), t)
	}

	return lhs
}

func (p *Processor) findRHS(node ast.Node) []string {
	var rhs []string

	if node == nil {
		return rhs
	}

	switch t := node.(type) {
	case *ast.BasicLit, *ast.SelectStmt, *ast.ChanType,
		*ast.LabeledStmt, *ast.DeclStmt, *ast.BranchStmt,
		*ast.TypeSpec, *ast.ArrayType, *ast.CaseClause,
		*ast.CommClause, *ast.KeyValueExpr, *ast.MapType,
		*ast.FuncLit:
	// Nothing to add to RHS
	case *ast.Ident:
		return []string{t.Name}
	case *ast.SelectorExpr:
		// TODO: Should this be RHS?
		// t.X is needed for defer as of now and t.Sel needed for special
		// functions such as Lock()
		rhs = p.findRHS(t.X)
		rhs = append(rhs, p.findRHS(t.Sel)...)
	case *ast.AssignStmt:
		for _, v := range t.Rhs {
			rhs = append(rhs, p.findRHS(v)...)
		}
	case *ast.CallExpr:
		for _, v := range t.Args {
			rhs = append(rhs, p.findRHS(v)...)
		}

		rhs = append(rhs, p.findRHS(t.Fun)...)
	case *ast.CompositeLit:
		for _, v := range t.Elts {
			rhs = append(rhs, p.findRHS(v)...)
		}
	case *ast.IfStmt:
		rhs = append(rhs, p.findRHS(t.Cond)...)
		rhs = append(rhs, p.findRHS(t.Init)...)
	case *ast.BinaryExpr:
		return append(
			p.findRHS(t.X),
			p.findRHS(t.Y)...,
		)
	case *ast.TypeSwitchStmt:
		return p.findRHS(t.Assign)
	case *ast.ReturnStmt:
		for _, v := range t.Results {
			rhs = append(rhs, p.findRHS(v)...)
		}
	case *ast.BlockStmt:
		for _, v := range t.List {
			rhs = append(rhs, p.findRHS(v)...)
		}
	case *ast.SwitchStmt:
		return p.findRHS(t.Tag)
	case *ast.GoStmt:
		return p.findRHS(t.Call)
	case *ast.ForStmt:
		return p.findRHS(t.Cond)
	case *ast.DeferStmt:
		return p.findRHS(t.Call)
	case *ast.SendStmt:
		return p.findLHS(t.Value)
	case *ast.IndexExpr:
		rhs = append(rhs, p.findRHS(t.Index)...)
		rhs = append(rhs, p.findRHS(t.X)...)
	case *ast.SliceExpr:
		rhs = append(rhs, p.findRHS(t.X)...)
		rhs = append(rhs, p.findRHS(t.Low)...)
		rhs = append(rhs, p.findRHS(t.High)...)
	default:
		if x, ok := maybeX(t); ok {
			return p.findRHS(x)
		}

		p.addWarning("UNKNOWN RHS", t.Pos(), t)
	}

	return rhs
}

func (p *Processor) findBlockStmt(node ast.Node) []*ast.BlockStmt {
	var blocks []*ast.BlockStmt

	switch t := node.(type) {
	case *ast.AssignStmt:
		for _, x := range t.Rhs {
			blocks = append(blocks, p.findBlockStmt(x)...)
		}
	case *ast.CallExpr:
		blocks = append(blocks, p.findBlockStmt(t.Fun)...)
	case *ast.FuncLit:
		blocks = append(blocks, t.Body)
	case *ast.ExprStmt:
		blocks = append(blocks, p.findBlockStmt(t.X)...)
	case *ast.ReturnStmt:
		for _, x := range t.Results {
			blocks = append(blocks, p.findBlockStmt(x)...)
		}
	case *ast.DeferStmt:
		blocks = append(blocks, p.findBlockStmt(t.Call)...)
	case *ast.GoStmt:
		blocks = append(blocks, p.findBlockStmt(t.Call)...)
	}

	return blocks
}

// maybeX extracts the X field from an AST node and returns it with a true value
// if it exists. If the node doesn't have an X field nil and false is returned.
// Known fields with X that are handled:
// IndexExpr, ExprStmt, SelectorExpr, StarExpr, ParentExpr, TypeAssertExpr,
// RangeStmt, UnaryExpr, ParenExpr, SliceExpr, IncDecStmt.
func maybeX(node interface{}) (ast.Node, bool) {
	maybeHasX := reflect.Indirect(reflect.ValueOf(node)).FieldByName("X")
	if !maybeHasX.IsValid() {
		return nil, false
	}

	n, ok := maybeHasX.Interface().(ast.Node)
	if !ok {
		return nil, false
	}

	return n, true
}

func atLeastOneInListsMatch(listOne, listTwo []string) bool {
	sliceToMap := func(s []string) map[string]struct{} {
		m := map[string]struct{}{}

		for _, v := range s {
			m[v] = struct{}{}
		}

		return m
	}

	m1 := sliceToMap(listOne)
	m2 := sliceToMap(listTwo)

	for k1 := range m1 {
		if _, ok := m2[k1]; ok {
			return true
		}
	}

	for k2 := range m2 {
		if _, ok := m1[k2]; ok {
			return true
		}
	}

	return false
}

// findLeadingAndTrailingWhitespaces will find leading and trailing whitespaces
// in a node. The method takes comments in consideration which will make the
// parser more gentle.
// nolint: gocognit
func (p *Processor) findLeadingAndTrailingWhitespaces(ident *ast.Ident, stmt, nextStatement ast.Node) {
	var (
		allowedLinesBeforeFirstStatement = 1
		commentMap                       = ast.NewCommentMap(p.fileSet, stmt, p.file.Comments)
		blockStatements                  []ast.Stmt
		blockStartLine                   int
		blockEndLine                     int
		blockStartPos                    token.Pos
		blockEndPos                      token.Pos
	)

	// Depending on the block type, get the statements in the block and where
	// the block starts (and ends).
	switch t := stmt.(type) {
	case *ast.BlockStmt:
		blockStatements = t.List
		blockStartPos = t.Lbrace
		blockEndPos = t.Rbrace
	case *ast.CaseClause:
		blockStatements = t.Body
		blockStartPos = t.Colon
	case *ast.CommClause:
		blockStatements = t.Body
		blockStartPos = t.Colon
	default:
		p.addWarning("whitespace node type not implemented ", stmt.Pos(), stmt)

		return
	}

	// Ignore empty blocks even if they have newlines or just comments.
	if len(blockStatements) < 1 {
		return
	}

	blockStartLine = p.fileSet.Position(blockStartPos).Line
	blockEndLine = p.fileSet.Position(blockEndPos).Line

	// No whitespace possible if LBrace and RBrace is on the same line.
	if blockStartLine == blockEndLine {
		return
	}

	var (
		firstStatement = blockStatements[0]
		lastStatement  = blockStatements[len(blockStatements)-1]
	)

	// Get the comment related to the first statement, we do allow commends in
	// the beginning of a block before the first statement.
	if c, ok := commentMap[firstStatement]; ok {
		for _, commentGroup := range c {
			// If the comment group is on the same line as the block start
			// (LBrace) we should not consider it.
			if p.nodeStart(commentGroup) == blockStartLine {
				continue
			}

			// We only care about comments before our statement from the comment
			// map. As soon as we hit comments after our statement let's break
			// out!
			if commentGroup.Pos() > firstStatement.Pos() {
				break
			}

			// Support both /* multiline */ and //single line comments
			for _, c := range commentGroup.List {
				allowedLinesBeforeFirstStatement += len(strings.Split(c.Text, "\n"))
			}
		}
	}

	if p.nodeStart(firstStatement) != blockStartLine+allowedLinesBeforeFirstStatement {
		p.addError(
			blockStartPos,
			"block should not start with a whitespace",
		)
	}

	// If the blockEndLine is not 0 we're a regular block (not case).
	if blockEndLine != 0 {
		if p.config.AllowTrailingComment {
			if lastComment, ok := commentMap[lastStatement]; ok {
				var (
					lastCommentGroup = lastComment[len(lastComment)-1]
					lastCommentLine  = lastCommentGroup.List[len(lastCommentGroup.List)-1]
					countNewlines    = 0
				)

				countNewlines += len(strings.Split(lastCommentLine.Text, "\n"))

				// No newlines between trailing comments and end of block.
				if p.nodeStart(lastCommentLine)+countNewlines != blockEndLine-1 {
					return
				}
			}
		}

		if p.nodeEnd(lastStatement) != blockEndLine-1 && !isExampleFunc(ident) {
			p.addError(blockEndPos, "block should not end with a whitespace (or comment)")
		}

		return
	}

	// If we don't have any nextStatement the trailing whitespace will be
	// handled when parsing the switch. If we do have a next statement we can
	// see where it starts by getting it's colon position. We set the end of the
	// current case to the position of the next case.
	switch n := nextStatement.(type) {
	case *ast.CaseClause:
		blockEndPos = n.Case
	case *ast.CommClause:
		blockEndPos = n.Case
	default:
		// No more cases
		return
	}

	blockEndLine = p.fileSet.Position(blockEndPos).Line - 1

	var (
		blockSize                = blockEndLine - blockStartLine
		caseTrailingCommentLines int
	)

	// TODO: I don't know what comments are bound to in cases. For regular
	// blocks the last comment is bound to the last statement but for cases
	// they are bound to the case clause expression. This will however get us all
	// comments and depending on the case expression this gets tricky.
	//
	// To handle this I get the comment map from the current statement (the case
	// itself) and iterate through all groups and all comment within all groups.
	// I then get the comments after the last statement but before the next case
	// clause and just map each line of comment that way.
	for _, commentGroups := range commentMap {
		for _, commentGroup := range commentGroups {
			for _, comment := range commentGroup.List {
				commentLine := p.fileSet.Position(comment.Pos()).Line

				// Ignore comments before the last statement.
				if commentLine <= p.nodeStart(lastStatement) {
					continue
				}

				// Ignore comments after the end of this case.
				if commentLine > blockEndLine {
					continue
				}

				// This allows /* multiline */ comments with newlines as well
				// as regular (//) ones
				caseTrailingCommentLines += len(strings.Split(comment.Text, "\n"))
			}
		}
	}

	hasTrailingWhitespace := p.nodeEnd(lastStatement)+caseTrailingCommentLines != blockEndLine

	// If the force trailing limit is configured and we don't end with a newline.
	if p.config.CaseForceTrailingWhitespaceLimit > 0 && !hasTrailingWhitespace {
		// Check if the block size is too big to miss the newline.
		if blockSize >= p.config.CaseForceTrailingWhitespaceLimit {
			p.addError(lastStatement.Pos(), "case block should end with newline at this size")
		}
	}
}

func isExampleFunc(ident *ast.Ident) bool {
	return ident != nil && strings.HasPrefix(ident.Name, "Example")
}

func (p *Processor) nodeStart(node ast.Node) int {
	return p.fileSet.Position(node.Pos()).Line
}

func (p *Processor) nodeEnd(node ast.Node) int {
	return p.fileSet.Position(node.End()).Line
}

// Add an error for the file and line number for the current token.Pos with the
// given reason.
func (p *Processor) addError(pos token.Pos, reason string) {
	position := p.fileSet.Position(pos)

	p.result = append(p.result, Result{
		FileName:   position.Filename,
		LineNumber: position.Line,
		Position:   position,
		Reason:     reason,
	})
}

func (p *Processor) addWarning(w string, pos token.Pos, t interface{}) {
	position := p.fileSet.Position(pos)

	p.warnings = append(p.warnings,
		fmt.Sprintf("%s:%d: %s (%T)", position.Filename, position.Line, w, t),
	)
}
