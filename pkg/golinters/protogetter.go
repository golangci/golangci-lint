package golinters

import (
	"sync"

	"github.com/ghostiam/protogetter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func NewProtoGetter() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	a := protogetter.NewAnalyzer()
	a.Run = func(pass *analysis.Pass) (any, error) {
		pgIssues := protogetter.Run(pass, protogetter.GolangciLintMode)

		issues := make([]goanalysis.Issue, len(pgIssues))
		for i, issue := range pgIssues {
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
