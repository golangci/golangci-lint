package golinters

import (
	"fmt"
	"sync"

	"github.com/yeya24/promlinter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const promlinterName = "promlinter"

func NewPromlinter(settings *config.PromlinterSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	var promSettings promlinter.Setting
	if settings != nil {
		promSettings = promlinter.Setting{
			Strict:            settings.Strict,
			DisabledLintFuncs: settings.DisabledLinters,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: promlinterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues := runPromLinter(pass, promSettings)

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
		promlinterName,
		"Check Prometheus metrics naming via promlint",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPromLinter(pass *analysis.Pass, promSettings promlinter.Setting) []goanalysis.Issue {
	lintIssues := promlinter.RunLint(pass.Fset, pass.Files, promSettings)

	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
		issue := result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("Metric: %s Error: %s", i.Metric, i.Text),
			FromLinter: promlinterName,
		}

		issues[k] = goanalysis.NewIssue(&issue, pass)
	}

	return issues
}
