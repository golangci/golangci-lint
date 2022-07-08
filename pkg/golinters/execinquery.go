package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/lufeee/execinquery"
	"golang.org/x/tools/go/analysis"
)

func NewExecInQuery() *goanalysis.Linter {
	a := execinquery.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
