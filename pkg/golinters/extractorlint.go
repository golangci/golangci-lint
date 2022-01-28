package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/nozzle/extractorlint/pkg/analyzer"
)

func NewExtractorLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"extractorlint",
		"checks extractor handlers for lint issues",
		[]*analysis.Analyzer{analyzer.HandlerAnalyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
