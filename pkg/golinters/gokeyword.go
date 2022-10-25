package golinters

import (
	"go/ast"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const (
	goKeywordName        = "gokeyword"
	goKeywordErrorMsg    = "detected use of go keyword: %s"
	goKeywordDescription = "detects presence of the go keyword"
	defaultDetails       = "no details provided"
)

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {
	return goanalysis.NewLinter(
		goKeywordName,
		goKeywordDescription,
		[]*analysis.Analyzer{{
			Name: goKeywordName,
			Doc:  goKeywordDescription,
			Run: func(pass *analysis.Pass) (interface{}, error) {
				details := defaultDetails
				if cfg != nil && cfg.Details != "" {
					details = cfg.Details
				}

				i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
				if !ok {
					return nil, errors.New("analyzer is not type *inspector.Inspector")
				}

				i.Preorder([]ast.Node{(*ast.GoStmt)(nil)}, func(node ast.Node) {
					if _, ok := node.(*ast.GoStmt); ok {
						pass.Reportf(node.Pos(), goKeywordErrorMsg, details)
					}
				})
				return nil, nil
			},
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		}},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
