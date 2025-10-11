package unnecessaryerror

import (
	"github.com/sollniss/unnecessaryerror/unnecessaryerror"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(unnecessaryerror.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
