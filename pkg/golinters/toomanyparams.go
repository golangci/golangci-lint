package golinters

import (
	"fmt"
	"go/ast"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	toomanyparams = "toomanyparams"

	maxIncomingParams = 5
	maxOutgoingParams = 3
)

//nolint:dupl
func NewTooManyParams(settings *config.TooManyParamsSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: toomanyparams,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runTooManyParams(pass, settings)
			if err != nil {
				return nil, err
			}

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
		toomanyparams,
		"Reports too many func paramters and return values",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runTooManyParams(pass *analysis.Pass, settings *config.TooManyParamsSettings) ([]goanalysis.Issue, error) {
	var issues []goanalysis.Issue

	for _, f := range pass.Files {
		newIssues := runFuncTooManyParams(pass, f, settings)
		issues = append(issues, newIssues...)
	}

	return issues, nil
}

func runFuncTooManyParams(pass *analysis.Pass, f *ast.File, settings *config.TooManyParamsSettings) []goanalysis.Issue {
	var issues []goanalysis.Issue
	for _, decl := range f.Decls {
		fnDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fnType := fnDecl.Type
		if fnType.Params != nil {
			num := 0
			for _, l := range fnType.Params.List {
				num += len(l.Names)
			}
			if num > settings.MaxParams {
				text := fmt.Sprintf("too many parameters (%d>%d) in func %s",
					num, settings.MaxParams, fnDecl.Name.Name)
				issues = append(issues, goanalysis.Issue{
					Issue: result.Issue{
						FromLinter: toomanyparams,
						Text:       text,
						Severity:   "warn",
						Pos:        pass.Fset.Position(fnType.Params.Pos()),
					},
					Pass: pass,
				})
			}
		}
		if fnType.Results != nil {
			if num := len(fnType.Results.List); num > settings.MaxReturns {
				text := fmt.Sprintf("too many return values (%d>%d) in func %s",
					num, settings.MaxReturns, fnDecl.Name.Name)
				issues = append(issues, goanalysis.Issue{
					Issue: result.Issue{
						FromLinter: toomanyparams,
						Text:       text,
						Severity:   "warn",
						Pos:        pass.Fset.Position(fnType.Results.Pos()),
					},
					Pass: pass,
				})
			}
		}
	}
	return issues
}
