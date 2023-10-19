package golinters

import (
	"github.com/catenacyber/perfsprint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewPerfSprint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"perfsprint",
		"Checks usages of `fmt.Sprintf` which have faster alternatives.",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
