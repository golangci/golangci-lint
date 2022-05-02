package golinters

import (
	"sync"

	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gofumptName = "gofumpt"

func NewGofumpt() *goanalysis.Linter {
	var (
		mu        sync.Mutex
		resIssues []goanalysis.Issue
	)

	analyzer := &analysis.Analyzer{
		Name: gofumptName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gofumptName,
		"Gofumpt checks whether code was gofumpt-ed.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		settings := lintCtx.Settings().Gofumpt

		options := format.Options{
			LangVersion: getLangVersion(settings),
			ModulePath:  settings.ModulePath,
			ExtraRules:  settings.ExtraRules,
		}

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			cb := func(_ string, src []byte) ([]byte, error) {
				return format.Source(src, options)
			}

			issues, err := runFormatAndDiffLinter(pass, lintCtx, gofumptName, cb)
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

func getLangVersion(settings config.GofumptSettings) string {
	if settings.LangVersion == "" {
		// TODO: defaults to "1.15", in the future (v2) must be set by using build.Default.ReleaseTags like staticcheck.
		return "1.15"
	}
	return settings.LangVersion
}
