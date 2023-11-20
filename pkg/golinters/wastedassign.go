package golinters

import (
	"github.com/sanposhiho/wastedassign/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewWastedAssign() *goanalysis.Linter {
	a := wastedassign.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Finds wasted assignment statements",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
