package gochecknoglobals

import (
	"4d63.com/gochecknoglobals/checknoglobals"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := checknoglobals.Analyzer()

	return goanalysis.NewLinter(
		a.Name,
		"Check that no global variables exist.",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
