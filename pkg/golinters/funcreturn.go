package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/hsnks100/funcreturn/pkg/funcreturn"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const linterName = "funcreturn"

func NewFuncReturn() *goanalysis.Linter {
	return goanalysis.NewLinter(
		linterName,
		"Checks whether there is a newline between functions",
		[]*analysis.Analyzer{
			funcreturn.Analyzer,
		},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
