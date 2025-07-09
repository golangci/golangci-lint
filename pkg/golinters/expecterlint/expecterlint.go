package expecterlint

import (
	"github.com/d0ubletr0uble/expecterlint"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(expecterlint.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
