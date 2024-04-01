package nestif

import (
	"sort"
	"sync"

	"github.com/nakabonne/nestif"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "nestif"

func New(settings *config.NestifSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runNestIf(pass, settings)

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
		name,
		"Reports deeply nested if statements",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNestIf(pass *analysis.Pass, settings *config.NestifSettings) []goanalysis.Issue {
	checker := &nestif.Checker{
		MinComplexity: settings.MinComplexity,
	}

	var lintIssues []nestif.Issue
	for _, f := range pass.Files {
		lintIssues = append(lintIssues, checker.Check(f, pass.Fset)...)
	}

	if len(lintIssues) == 0 {
		return nil
	}

	sort.SliceStable(lintIssues, func(i, j int) bool {
		return lintIssues[i].Complexity > lintIssues[j].Complexity
	})

	issues := make([]goanalysis.Issue, 0, len(lintIssues))
	for _, i := range lintIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: name,
		}, pass))
	}

	return issues
}
