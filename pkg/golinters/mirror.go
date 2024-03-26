package golinters

import (
	"sync"

	"github.com/butuzov/mirror"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func NewMirror() *goanalysis.Linter {
	var (
		mu     sync.Mutex
		issues []goanalysis.Issue
	)

	a := mirror.NewAnalyzer()
	a.Run = func(pass *analysis.Pass) (any, error) {
		// mirror only lints test files if the `--with-tests` flag is passed,
		// so we pass the `with-tests` flag as true to the analyzer before running it.
		// This can be turned off by using the regular golangci-lint flags such as `--tests` or `--skip-files`
		// or can be disabled per linter via exclude rules.
		// (see https://github.com/golangci/golangci-lint/issues/2527#issuecomment-1023707262)
		violations := mirror.Run(pass, true)

		if len(violations) == 0 {
			return nil, nil
		}

		for index := range violations {
			i := violations[index].Issue(pass.Fset)

			issue := result.Issue{
				FromLinter: a.Name,
				Text:       i.Message,
				Pos:        i.Start,
			}

			if i.InlineFix != "" {
				issue.Replacement = &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  i.Start.Column - 1,
						Length:    len(i.Original),
						NewString: i.InlineFix,
					},
				}
			}

			mu.Lock()
			issues = append(issues, goanalysis.NewIssue(&issue, pass))
			mu.Unlock()
		}

		return nil, nil
	}

	analyzer := goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return issues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)

	return analyzer
}
