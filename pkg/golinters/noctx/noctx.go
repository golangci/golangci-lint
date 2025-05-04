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
		"Detects function and method with missing usage of context.Context",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
