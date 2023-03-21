package golinters

import (
	"sync"

	"github.com/4meepo/tagalign"
	"github.com/leonklingele/grouper/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const tagalignName = "tagalign"

func NewTagAlign(settings *config.TagAlignSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	a := &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check if struct tags are well aligned",
		Run: func(p *analysis.Pass) (any, error) {
			var options []tagalign.Option
			options = append(options, tagalign.WithMode(tagalign.GolangciLintMode))
			if settings.AutoSort {
				if len(settings.FixedOrder) > 0 {
					options = append(options, tagalign.WithAutoSort(settings.FixedOrder...))
				} else {
					options = append(options, tagalign.WithAutoSort())
				}
			}

			tagalignIssues := tagalign.Run(p, options...)

			issues := make([]goanalysis.Issue, len(tagalignIssues))
			for i, issue := range tagalignIssues {
				replacement := result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  issue.InlineFix.StartCol,
						Length:    issue.InlineFix.Length,
						NewString: issue.InlineFix.NewString,
					},
				}
				issues[i] = goanalysis.NewIssue(&result.Issue{
					FromLinter:  tagalignName,
					Pos:         issue.Pos,
					Text:        issue.Message,
					Replacement: &replacement,
				}, p)
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
		tagalignName,
		analyzer.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
