package golinters

import (
	"golang.org/x/tools/go/analysis"
	"tildegit.org/indigo/explicitreturn/pkg/analyzer"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExplicitReturn() *goanalysis.Linter {
	a := analyzer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
