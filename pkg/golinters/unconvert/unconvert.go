package unconvert

import (
	"sync"

	"github.com/golangci/unconvert"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const linterName = "unconvert"

func New(settings *config.UnconvertSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runUnconvert(pass, settings)

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
		linterName,
		"Remove unnecessary type conversions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnconvert(pass *analysis.Pass, settings *config.UnconvertSettings) []goanalysis.Issue {
	positions := unconvert.Run(pass, settings.FastMath, settings.Safe)

	var issues []goanalysis.Issue
	for _, position := range positions {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        position,
			Text:       "unnecessary conversion",
			FromLinter: linterName,
		}, pass))
	}

	return issues
}
