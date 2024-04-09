package gocyclo

import (
	"fmt"
	"sync"

	"github.com/fzipp/gocyclo"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "gocyclo"

func New(settings *config.GoCycloSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runGoCyclo(pass, settings)

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
		"Computes and checks the cyclomatic complexity of functions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoCyclo(pass *analysis.Pass, settings *config.GoCycloSettings) []goanalysis.Issue {
	var stats gocyclo.Stats
	for _, f := range pass.Files {
		stats = gocyclo.AnalyzeASTFile(f, pass.Fset, stats)
	}
	if len(stats) == 0 {
		return nil
	}

	stats = stats.SortAndFilter(-1, settings.MinComplexity)

	issues := make([]goanalysis.Issue, 0, len(stats))

	for _, s := range stats {
		text := fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
			s.Complexity, internal.FormatCode(s.FuncName, nil), settings.MinComplexity)

		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        s.Pos,
			Text:       text,
			FromLinter: name,
		}, pass))
	}

	return issues
}
