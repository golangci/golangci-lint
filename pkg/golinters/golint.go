package golinters

import (
	"fmt"
	"sync"

	lintAPI "github.com/golangci/lint-1"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const golintName = "golint"

//nolint:dupl
func NewGolint(settings *config.GoLintSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: golintName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runGoLint(pass, settings)
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
		golintName,
		"Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runGoLint(pass *analysis.Pass, settings *config.GoLintSettings) ([]goanalysis.Issue, error) {
	l := new(lintAPI.Linter)

	ps, err := l.LintPkg(pass.Files, pass.Fset, pass.Pkg, pass.TypesInfo)
	if err != nil {
		return nil, fmt.Errorf("can't lint %d files: %w", len(pass.Files), err)
	}

	if len(ps) == 0 {
		return nil, nil
	}

	lintIssues := make([]*result.Issue, 0, len(ps)) // This is worst case
	for idx := range ps {
		if ps[idx].Confidence >= settings.MinConfidence {
			lintIssues = append(lintIssues, &result.Issue{
				Pos:        ps[idx].Position,
				Text:       ps[idx].Text,
				FromLinter: golintName,
			})
			// TODO: use p.Link and p.Category
		}
	}

	issues := make([]goanalysis.Issue, 0, len(lintIssues))
	for _, issue := range lintIssues {
		issues = append(issues, goanalysis.NewIssue(issue, pass))
	}

	return issues, nil
}
