package bodyclose

import (
	"github.com/timakin/bodyclose/passes/bodyclose"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(bodyclose.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
