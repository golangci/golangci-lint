package golinters

import (
	"github.com/esimonov/ifshort/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewIfshort() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"ifshort",
		"Checks that your code uses short syntax for if-statements whenever possible",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
