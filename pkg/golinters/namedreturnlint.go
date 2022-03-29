package golinters

import (
	"github.com/firefart/namedreturnlint/analyzer"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"golang.org/x/tools/go/analysis"
)

func NewNamedReturnLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"namedreturnlint",
		"Checks for named returns",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
