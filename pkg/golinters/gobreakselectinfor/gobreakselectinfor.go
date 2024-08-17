package gobreakselectinfor

import (
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const linterName = "gobreakselectinfor"

func New() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			fileIssues := run(pass)
			res := make([]goanalysis.Issue, 0, len(fileIssues))
			for i := range fileIssues {
				res = append(res, goanalysis.NewIssue(&fileIssues[i], pass))
			}
			if len(res) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks that break statement inside select statement inside for loop",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(pass *analysis.Pass) []result.Issue {
	var res []result.Issue

	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		ast.Inspect(funcDecl.Body, func(stmt ast.Node) bool {
			if forStmt, ok := stmt.(*ast.ForStmt); ok {
				ast.Inspect(forStmt.Body, func(stmt ast.Node) bool {
					if selStmt, ok := stmt.(*ast.SelectStmt); ok {
						ast.Inspect(selStmt.Body, func(stmt ast.Node) bool {
							if brkStmt, ok := stmt.(*ast.BranchStmt); ok && brkStmt.Tok == token.BREAK {
								pass.Reportf(stmt.Pos(), "break statement inside select statement inside for loop")
								res = append(res, result.Issue{
									Pos:        pass.Fset.Position(stmt.Pos()),
									Text:       "break statement inside select statement inside for loop",
									FromLinter: linterName,
								})
								return true
							}
							return true
						})
					}
					return true
				})
			}
			return true
		})

		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}

	return res
}
