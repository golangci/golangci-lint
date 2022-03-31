package golinters

import (
	"github.com/firefart/namedreturnlint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNamedReturnLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"namedreturnlint",
		"Checks for named returns",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
