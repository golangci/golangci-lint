package golinters

import (
	"sync"

	goimportsAPI "github.com/golangci/gofmt/goimports"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const goimportsName = "goimports"

func NewGoimports(settings *config.GoImportsSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goimportsName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		goimportsName,
		"Check import statements are formatted according to the 'goimport' command. "+
			"Reformat imports in autofix mode.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		imports.LocalPrefix = settings.LocalPrefixes

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runGoImports(lintCtx, pass)
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
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoImports(lintCtx *linter.Context, pass *analysis.Pass) ([]goanalysis.Issue, error) {
	fileNames := getFileNames(pass)

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		diff, err := goimportsAPI.Run(f)
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := extractIssuesFromPatch(string(diff), lintCtx, goimportsName)
		if err != nil {
			return nil, errors.Wrapf(err, "can't extract issues from gofmt diff output %q", string(diff))
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}
