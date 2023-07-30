package golinters

import (
	"github.com/HarryTennent/noimportsasvar/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewNoImportsAsVar() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"noimportsasvar",
		"Checks that a file's imports are not used as variable names.",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
