package nilerr

import (
	"github.com/gostaticanalysis/nilerr"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := nilerr.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Finds the code that returns nil even if it checks that the error is not nil.",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
