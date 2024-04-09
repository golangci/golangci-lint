package gochecknoinits

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "gochecknoinits"

func New() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			var res []goanalysis.Issue
			for _, file := range pass.Files {
				fileIssues := checkFileForInits(file, pass.Fset)
				for i := range fileIssues {
					res = append(res, goanalysis.NewIssue(&fileIssues[i], pass))
				}
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
		name,
		"Checks that no init functions are present in Go code",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func checkFileForInits(f *ast.File, fset *token.FileSet) []result.Issue {
	var res []result.Issue
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := funcDecl.Name.Name
		if name == "init" && funcDecl.Recv.NumFields() == 0 {
			res = append(res, result.Issue{
				Pos:        fset.Position(funcDecl.Pos()),
				Text:       fmt.Sprintf("don't use %s function", internal.FormatCode(name, nil)),
				FromLinter: name,
			})
		}
	}

	return res
}
