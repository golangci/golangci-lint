package goimports

import (
	"fmt"
	"sync"

	goimportsAPI "github.com/golangci/gofmt/goimports"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const name = "goimports"

func New(settings *config.GoImportsSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		name,
		"Check import statements are formatted according to the 'goimport' command. "+
			"Reformat imports in autofix mode.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		imports.LocalPrefix = settings.LocalPrefixes

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
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
	fileNames := internal.GetFileNames(pass)

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		diff, err := goimportsAPI.Run(f)
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := internal.ExtractIssuesFromPatch(string(diff), lintCtx, name, getIssuedTextGoImports)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %w", string(diff), err)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

func getIssuedTextGoImports(settings *config.LintersSettings) string {
	text := "File is not `goimports`-ed"

	if settings.Goimports.LocalPrefixes != "" {
		text += " with -local " + settings.Goimports.LocalPrefixes
	}

	return text
}
