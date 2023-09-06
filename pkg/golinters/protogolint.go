package golinters

import (
	"sync"

	"github.com/ghostiam/protogolint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func NewProtoGoLint() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	a := protogolint.NewAnalyzer()
	a.Run = func(pass *analysis.Pass) (any, error) {
		taIssues := protogolint.Run(pass, protogolint.GolangciLintMode)

		issues := make([]goanalysis.Issue, len(taIssues))
		for i, issue := range taIssues {
			report := &result.Issue{
				FromLinter: a.Name,
				Pos:        issue.Pos,
				Text:       issue.Message,
				Replacement: &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  issue.InlineFix.StartCol,
						Length:    issue.InlineFix.Length,
						NewString: issue.InlineFix.NewString,
					},
				},
			}

			issues[i] = goanalysis.NewIssue(report, pass)
		}

		if len(issues) == 0 {
			return nil, nil
		}

		mu.Lock()
		resIssues = append(resIssues, issues...)
		mu.Unlock()

		return nil, nil
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
