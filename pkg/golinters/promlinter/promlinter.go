package promlinter

import (
	"fmt"

	"github.com/yeya24/promlinter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

const linterName = "promlinter"

func New(settings *config.PromlinterSettings) *goanalysis.Linter {
	b := goanalysis.NewThreadSafeLinterBuilder()

	var promSettings promlinter.Setting
	if settings != nil {
		promSettings = promlinter.Setting{
			Strict:            settings.Strict,
			DisabledLintFuncs: settings.DisabledLinters,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: linterName,
			Doc:  "Check Prometheus metrics naming via promlint",
			Run: func(pass *analysis.Pass) (any, error) {
				b.Add(runPromLinter(pass, promSettings)...)
				return nil, nil
			},
		}).
		WithIssuesReporter(b.Reporter()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPromLinter(pass *analysis.Pass, promSettings promlinter.Setting) []*goanalysis.Issue {
	lintIssues := promlinter.RunLint(pass.Fset, pass.Files, promSettings)

	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]*goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
		issue := result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("Metric: %s Error: %s", i.Metric, i.Text),
			FromLinter: linterName,
		}

		issues[k] = goanalysis.NewIssue(&issue, pass)
	}

	return issues
}
