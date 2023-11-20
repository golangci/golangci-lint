package golinters

import (
	"github.com/tdakkota/asciicheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewAsciicheck() *goanalysis.Linter {
	a := asciicheck.NewAnalyzer()
	return goanalysis.NewLinter(
		a.Name,
		"Simple linter to check that your code does not contain non-ASCII identifiers",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
