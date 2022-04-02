package golinters

import (
	"github.com/firefart/nonamedreturnlint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoNamedReturnLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"nonamedreturnlint",
		"Reports all named returns",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
