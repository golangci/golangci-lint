package golinters

import (
	"fmt"
	"go/ast"
	"sync"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const NoLintLintName = "nolintlint"

//nolint:dupl
func NewNoLintLint(settings *config.NoLintLintSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: NoLintLintName,
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
		NoLintLintName,
		"Reports ill-formed or insufficient nolint directives",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runNoLintLint(pass *analysis.Pass, settings *config.NoLintLintSettings) ([]goanalysis.Issue, error) {
	var needs nolintlint.Needs
	if settings.RequireExplanation {
		needs |= nolintlint.NeedsExplanation
	}
	if settings.RequireSpecific {
		needs |= nolintlint.NeedsSpecific
	}
	if !settings.AllowUnused {
		needs |= nolintlint.NeedsUnused
	}

	lnt, err := nolintlint.NewLinter(needs, settings.AllowNoExplanation)
	if err != nil {
		return nil, err
	}

	nodes := make([]ast.Node, 0, len(pass.Files))
	for _, n := range pass.Files {
		nodes = append(nodes, n)
	}

	lintIssues, err := lnt.Run(pass.Fset, nodes...)
	if err != nil {
		return nil, fmt.Errorf("linter failed to run: %s", err)
	}

	var issues []goanalysis.Issue

	for _, i := range lintIssues {
		expectNoLint := false
		var expectedNolintLinter string
		if ii, ok := i.(nolintlint.UnusedCandidate); ok {
			expectedNolintLinter = ii.ExpectedLinter
			expectNoLint = true
		}

		issue := &result.Issue{
			FromLinter:           NoLintLintName,
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
