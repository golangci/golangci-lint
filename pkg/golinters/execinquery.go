package golinters

import (
	"github.com/lufeee/execinquery"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExecInQuery() *goanalysis.Linter {
	const linterName = "execinquery"

	a := execinquery.Analyzer
	a.Name = linterName // TODO the name must change inside the linter.

	return goanalysis.NewLinter(
		linterName,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
