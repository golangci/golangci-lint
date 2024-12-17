package gochecknoinits

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

const linterName = "gochecknoinits"

func New() *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name:     linterName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks that no init functions are present in Go code",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(decl ast.Node) {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			return
		}

		fnName := funcDecl.Name.Name
		if fnName == "init" && funcDecl.Recv.NumFields() == 0 {
			pass.Reportf(funcDecl.Pos(), "don't use %s function", internal.FormatCode(fnName, nil))
		}
	})

	return nil, nil
}
