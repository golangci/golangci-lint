package golinters

import (
	"github.com/gardener/gardener/hack/tools/gomegacheck/pkg/gomegacheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGomegaCheck() *goanalysis.Linter {
	a := gomegacheck.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)
}
