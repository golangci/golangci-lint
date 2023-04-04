package golinters

import (
	"github.com/sanposhiho/wastedassign/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewWastedAssign() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"wastedassign",
		"wastedassign finds wasted assignment statements.",
		[]*analysis.Analyzer{wastedassign.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
