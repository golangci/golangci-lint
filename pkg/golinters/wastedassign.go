package golinters

import (
	"github.com/sanposhiho/wastedassign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewWastedAssign() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		wastedassign.Analyzer,
	}

	return goanalysis.NewLinter(
		"wastedassign",
		"wastedassign finds wasted assignment statements.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
