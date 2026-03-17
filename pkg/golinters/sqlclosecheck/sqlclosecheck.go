package sqlclosecheck

import (
	"github.com/golangci/sqlclosecheck/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewDeferOnlyAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
