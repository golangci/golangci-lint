package golinters

import (
	"sync"

	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func NewTagAlign(settings *config.TagAlignSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	options := []tagalign.Option{tagalign.WithMode(tagalign.GolangciLintMode)}

	if settings != nil {
		options = append(options, tagalign.WithAlign(settings.Align))

		if settings.Sort || len(settings.Order) > 0 {
			options = append(options, tagalign.WithSort(settings.Order...))
		}
	}

	analyzer := tagalign.NewAnalyzer(options...)
	analyzer.Run = func(pass *analysis.Pass) (any, error) {
		taIssues := tagalign.Run(pass, options...)

		issues := make([]goanalysis.Issue, len(taIssues))
		for i, issue := range taIssues {
			report := &result.Issue{
				FromLinter: analyzer.Name,
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
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
