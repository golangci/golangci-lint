package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const dogsledName = "dogsled"

func NewDogsled(settings *config.DogsledSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: dogsledName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runDogsled(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		dogsledName,
		"Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runDogsled(pass *analysis.Pass, settings *config.DogsledSettings) []goanalysis.Issue {
	var reports []goanalysis.Issue
	for _, f := range pass.Files {
		v := &returnsVisitor{
			maxBlanks: settings.MaxBlankIdentifiers,
			f:         pass.Fset,
		}

		ast.Walk(v, f)

		for i := range v.issues {
			reports = append(reports, goanalysis.NewIssue(&v.issues[i], pass))
		}
	}

	return reports
}

type returnsVisitor struct {
	f         *token.FileSet
	maxBlanks int
	issues    []result.Issue
}

func (v *returnsVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}
	if funcDecl.Body == nil {
		return v
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

		if numBlank > v.maxBlanks {
			v.issues = append(v.issues, result.Issue{
				FromLinter: dogsledName,
				Text:       fmt.Sprintf("declaration has %v blank identifiers", numBlank),
				Pos:        v.f.Position(assgnStmt.Pos()),
			})
		}
	}
	return v
}
