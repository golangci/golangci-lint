package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"
)

func NewNoctx() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"noctx",
		"noctx finds sending http request without context.Context",
		[]*analysis.Analyzer{noctx.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
