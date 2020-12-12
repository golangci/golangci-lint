package golinters

import (
	"go/ast"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
	"github.com/anduril/golangci-lint/pkg/lint/linter"
	"github.com/anduril/golangci-lint/pkg/result"
)

const tickerInLoopName = "tickerinloop"

func NewTickerInLoop() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: tickerInLoopName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		tickerInLoopName,
		"Tool for detecting tickers created in loops",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues := runOnFile(pass)

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runOnFile(pass *analysis.Pass) []goanalysis.Issue {
	var res []goanalysis.Issue
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if forLoop, ok := n.(*ast.ForStmt); ok {
				for _, line := range forLoop.Body.List {
					assign, ok := line.(*ast.AssignStmt)
					if !ok {
						continue
					}
					for _, expr := range assign.Rhs {
						call, ok := expr.(*ast.CallExpr)
						if !ok {
							continue
						}
						selector, ok := call.Fun.(*ast.SelectorExpr)
						if !ok {
							continue
						}
						if selector.Sel.Name == "NewTicker" {
							res = append(res, goanalysis.NewIssue(&result.Issue{
								Pos:        pass.Fset.Position(selector.Sel.Pos()),
								Text:       "ticker found in loop",
								FromLinter: tickerInLoopName,
							}, pass))
						}
					}
				}
			}
			return true
		})
	}
	return res
}
