package noctx

import (
	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := noctx.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Finds sending http request without context.Context",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
