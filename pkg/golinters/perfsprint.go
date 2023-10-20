package golinters

import (
	"github.com/catenacyber/perfsprint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewPerfSprint() *goanalysis.Linter {
	a := analyzer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
