package golinters

import (
	"fmt"
	"sort"
	"sync"

	"github.com/uudashr/gocognit"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const gocognitName = "gocognit"

func NewGocognit(settings *config.GocognitSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runGocognit(pass, settings)

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
		gocognitName,
		"Computes and checks the cognitive complexity of functions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGocognit(pass *analysis.Pass, settings *config.GocognitSettings) []goanalysis.Issue {
	var stats []gocognit.Stat
	for _, f := range pass.Files {
		stats = gocognit.ComplexityStats(f, pass.Fset, stats)
	}
	if len(stats) == 0 {
		return nil
	}

	sort.SliceStable(stats, func(i, j int) bool {
		return stats[i].Complexity > stats[j].Complexity
	})

	issues := make([]goanalysis.Issue, 0, len(stats))
	for _, s := range stats {
		if s.Complexity <= settings.MinComplexity {
			break // Break as the stats is already sorted from greatest to least
		}

		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos: s.Pos,
			Text: fmt.Sprintf("cognitive complexity %d of func %s is high (> %d)",
				s.Complexity, internal.FormatCode(s.FuncName, nil), settings.MinComplexity),
			FromLinter: gocognitName,
		}, pass))
	}

	return issues
}
