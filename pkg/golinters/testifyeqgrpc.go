// nolint:dupl
package golinters

import (
	"go/ast"
	"go/types"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
	"github.com/anduril/golangci-lint/pkg/lint/linter"
	"github.com/anduril/golangci-lint/pkg/result"
)

const (
	testifyAssertEqualGrpc    = "testifyeqgrpc"
	testifyAssertEqualGrpcMsg = "call to assert.Equal made error type returned from " +
		"'google.golang.org/grpc/status': Use assert.EqualError or assert.Nil instead."
	grpcStatusPkg = "google.golang.org/grpc/status"
	stackSize     = 32
)

func NewTestifyAssertEqualGrpc() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     testifyAssertEqualGrpc,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      goanalysis.DummyRun,
	}
	return goanalysis.NewLinter(
		testifyAssertEqualGrpc,
		"Finds places in which assert.Equal with an error type returned from google.golang.org",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues := runTestifyEqGrpcStatus(pass)

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runTestifyEqGrpcStatus(pass *analysis.Pass) []goanalysis.Issue {
	var allIssues []goanalysis.Issue
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncLit)(nil),
		(*ast.FuncDecl)(nil),
	}
	i.Preorder(nodeFilter, func(n ast.Node) {
		allIssues = append(allIssues, runFunc(pass, n)...)
	})
	return allIssues
}

func runFunc(pass *analysis.Pass, n ast.Node) []goanalysis.Issue { // nolint:gocyclo
	var issues []goanalysis.Issue

	var statusVars []*ast.Ident
	stack := make([]ast.Node, 0, stackSize)
	ast.Inspect(n, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			if len(stack) > 0 {
				return false // don't stray into nested functions
			}
		case nil:
			stack = stack[:len(stack)-1] // pop
			return true
		}
		stack = append(stack, n) // push

		if !isStatusError(pass.TypesInfo, n) || !isCall(stack[len(stack)-2]) {
			return true
		}

		stmt := stack[len(stack)-3]

		switch stmt := stmt.(type) {
		case *ast.CallExpr:
			var id *ast.Ident
			switch fun := stmt.Fun.(type) {
			case *ast.Ident:
				id = fun
			case *ast.SelectorExpr:
				id = fun.Sel
			}
			use := pass.TypesInfo.Uses[id]
			if use.Pkg().Path() == "github.com/stretchr/testify/assert" && use.Id() == "Equal" {
				issue := goanalysis.NewIssue(&result.Issue{
					FromLinter: testifyAssertEqualGrpc,
					Pos:        pass.Fset.Position(stmt.Pos()),
					Text:       testifyAssertEqualGrpcMsg,
				}, pass)
				issues = append(issues, issue)
			}
		case *ast.AssignStmt:
			id, ok := stmt.Lhs[0].(*ast.Ident)
			if ok {
				statusVars = append(statusVars, id)
			}
		}
		return true
	})

	ast.Inspect(n, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return true // not a static call
		} else if fn.FullName() != testifyAssertEqualMethod {
			return true
		} else if len(call.Args) < 3 { // nolint:gomnd
			return true
		}

		for _, statusVar := range statusVars {
			bad := false
			if a1ident, ok := call.Args[1].(*ast.Ident); ok {
				if a1ident.Obj == statusVar.Obj {
					bad = true
				}
			}
			if a2ident, ok := call.Args[2].(*ast.Ident); ok {
				if a2ident.Obj == statusVar.Obj {
					bad = true
				}
			}
			if bad {
				issue := goanalysis.NewIssue(&result.Issue{
					FromLinter: testifyAssertEqualGrpc,
					Pos:        pass.Fset.Position(call.Pos()),
					Text:       testifyAssertEqualGrpcMsg,
				}, pass)
				issues = append(issues, issue)
			}
		}
		return true
	})
	return issues
}

// isStatusError reports whether n is one of the qualified identifiers status.Error{,f}.
func isStatusError(info *types.Info, n ast.Node) bool {
	sel, ok := n.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	switch sel.Sel.Name {
	case "New", "Newf", "Error", "Errorf", "ErrorProto", "FromProto":
	default:
		return false
	}
	if x, ok := sel.X.(*ast.Ident); ok {
		if pkgname, ok := info.Uses[x].(*types.PkgName); ok {
			return pkgname.Imported().Path() == grpcStatusPkg
		}
		// Import failed, so we can't check package path.
		// Just check the local package name (heuristic).
		return x.Name == grpcStatusPkg
	}
	return false
}

func isCall(n ast.Node) bool { _, ok := n.(*ast.CallExpr); return ok }
