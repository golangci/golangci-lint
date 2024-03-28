package nolintlint

import (
	"fmt"
	"go/ast"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const Name = "nolintlint"

func New(settings *config.NoLintLintSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: Name,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runNoLintLint(pass, settings)
			if err != nil {
				return nil, err
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
		Name,
		"Reports ill-formed or insufficient nolint directives",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNoLintLint(pass *analysis.Pass, settings *config.NoLintLintSettings) ([]goanalysis.Issue, error) {
	var needs internal.Needs
	if settings.RequireExplanation {
		needs |= internal.NeedsExplanation
	}
	if settings.RequireSpecific {
		needs |= internal.NeedsSpecific
	}
	if !settings.AllowUnused {
		needs |= internal.NeedsUnused
	}

	lnt, err := internal.NewLinter(needs, settings.AllowNoExplanation)
	if err != nil {
		return nil, err
	}

	nodes := make([]ast.Node, 0, len(pass.Files))
	for _, n := range pass.Files {
		nodes = append(nodes, n)
	}

	lintIssues, err := lnt.Run(pass.Fset, nodes...)
	if err != nil {
		return nil, fmt.Errorf("linter failed to run: %w", err)
	}

	var issues []goanalysis.Issue

	for _, i := range lintIssues {
		expectNoLint := false
		var expectedNolintLinter string
		if ii, ok := i.(internal.UnusedCandidate); ok {
			expectedNolintLinter = ii.ExpectedLinter
			expectNoLint = true
		}

		issue := &result.Issue{
			FromLinter:           Name,
			Text:                 i.Details(),
			Pos:                  i.Position(),
			ExpectNoLint:         expectNoLint,
			ExpectedNoLintLinter: expectedNolintLinter,
			Replacement:          i.Replacement(),
		}

		issues = append(issues, goanalysis.NewIssue(issue, pass))
	}

	return issues, nil
}
