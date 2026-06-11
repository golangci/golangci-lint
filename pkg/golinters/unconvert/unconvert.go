package unconvert

import (
	"github.com/golangci/unconvert"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

const linterName = "unconvert"

func New(settings *config.UnconvertSettings) *goanalysis.Linter {
	b := goanalysis.NewThreadSafeLinterBuilder()

	unconvert.SetFastMath(settings.FastMath)
	unconvert.SetSafe(settings.Safe)

	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: linterName,
			Doc:  "Remove unnecessary type conversions",
			Run: func(pass *analysis.Pass) (any, error) {
				b.Add(runUnconvert(pass)...)
				return nil, nil
			},
		}).
		WithIssuesReporter(b.Reporter()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnconvert(pass *analysis.Pass) []*goanalysis.Issue {
	positions := unconvert.Run(pass)

	var issues []*goanalysis.Issue
	for _, position := range positions {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        position,
			Text:       "unnecessary conversion",
			FromLinter: linterName,
		}, pass))
	}

	return issues
}
