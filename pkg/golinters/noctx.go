package golinters

import (
	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewNoctx() *goanalysis.Linter {
	a := noctx.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Finds sending http request without context.Context",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
