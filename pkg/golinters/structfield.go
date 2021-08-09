package golinters

import (
	"fmt"
	"go/ast"
	"sync"

	"github.com/uudashr/structfield"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const structfieldName = "structfield"

func NewStructfield() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: structfieldName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}

	return goanalysis.NewLinter(
		structfieldName,
		"Find struct literals using non-labeled fields",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			inspect := inspector.New(pass.Files)

			nodeFilter := []ast.Node{
				(*ast.CompositeLit)(nil),
			}
			var res []goanalysis.Issue
			inspect.Preorder(nodeFilter, func(n ast.Node) {
				lit := n.(*ast.CompositeLit)
				ok, count := structfield.CountNonLabeledFields(lit)
				if !ok {
					return
				}

				limit := lintCtx.Settings().Structfield.Limit
				if count > limit {
					pos := pass.Fset.Position(lit.Pos())
					res = append(res, goanalysis.NewIssue(&result.Issue{
						Pos:        pos,
						Text:       fmt.Sprintf("found %d non-labeled fields on struct literal (> %d)", count, limit),
						FromLinter: structfieldName,
					}, pass))
				}
			})

			if len(res) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
