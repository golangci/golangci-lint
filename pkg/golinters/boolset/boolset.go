package boolset

import (
	"github.com/arturmelanchyk/boolset/boolset"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(boolset.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
