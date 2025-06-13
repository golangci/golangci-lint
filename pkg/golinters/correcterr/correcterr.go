package correcterr

import (
	"github.com/m-ocean-it/correcterr"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(correcterr.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
