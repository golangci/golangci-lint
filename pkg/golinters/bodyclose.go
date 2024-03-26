package golinters

import (
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewBodyclose() *goanalysis.Linter {
	a := bodyclose.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"checks whether HTTP response body is closed successfully",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
