package varasconst

import (
	"github.com/MeenaAlfons/varasconst/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
