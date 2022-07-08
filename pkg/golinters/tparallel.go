package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/moricho/tparallel"
	"golang.org/x/tools/go/analysis"
)

func NewTparallel() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"tparallel",
		"tparallel detects inappropriate usage of t.Parallel() method in your Go test codes",
		[]*analysis.Analyzer{tparallel.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
