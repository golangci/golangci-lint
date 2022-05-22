package golinters

import (
	"sync"

	unconvertAPI "github.com/golangci/unconvert"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const unconvertName = "unconvert"

//nolint:dupl
func NewUnconvert() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: unconvertName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues := runUnconvert(pass)

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
		unconvertName,
		"Remove unnecessary type conversions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnconvert(pass *analysis.Pass) []goanalysis.Issue {
	prog := goanalysis.MakeFakeLoaderProgram(pass)

	positions := unconvertAPI.Run(prog)
	if len(positions) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, 0, len(positions))
	for _, pos := range positions {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        pos,
			Text:       "unnecessary conversion",
			FromLinter: unconvertName,
		}, pass))
	}

	return issues
}
