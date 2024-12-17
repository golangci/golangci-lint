package dogsled

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

const linterName = "dogsled"

func New(settings *config.DogsledSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, settings.MaxBlankIdentifiers)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass, maxBlanks int) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		if funcDecl.Body == nil {
			return
		}

		for _, expr := range funcDecl.Body.List {
			assgnStmt, ok := expr.(*ast.AssignStmt)
			if !ok {
				continue
			}

			numBlank := 0
			for _, left := range assgnStmt.Lhs {
				ident, ok := left.(*ast.Ident)
				if !ok {
					continue
				}
				if ident.Name == "_" {
					numBlank++
				}
			}

			if numBlank > maxBlanks {
				pass.Reportf(assgnStmt.Pos(), "declaration has %v blank identifiers", numBlank)
			}
		}
	})

	return nil, nil
}
