package nolintlint

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
	nolintlint "github.com/golangci/golangci-lint/v2/pkg/golinters/nolintlint/internal"
)

const LinterName = nolintlint.LinterName

func New(settings *config.NoLintLintSettings) *goanalysis.Linter {
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
		internal.LinterLogger.Fatalf("%s: create analyzer: %v", nolintlint.LinterName, err)
	}

	b := goanalysis.NewThreadSafeLinterBuilder()

	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: nolintlint.LinterName,
			Doc:  "Reports ill-formed or insufficient nolint directives",
			Run: func(pass *analysis.Pass) (any, error) {
				issues, err := lnt.Run(pass)
				if err != nil {
					return nil, fmt.Errorf("linter failed to run: %w", err)
				}

				b.Add(issues...)
				return nil, nil
			},
		}).
		WithIssuesReporter(b.Reporter()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
