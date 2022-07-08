package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/charithe/durationcheck"
	"golang.org/x/tools/go/analysis"
)

func NewDurationCheck() *goanalysis.Linter {
	a := durationcheck.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
