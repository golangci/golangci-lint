package golinters

import (
	"github.com/bedakb/nomainreturn"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoMainReturn() *goanalysis.Linter {
	a := nomainreturn.NewAnalyzer(nomainreturn.DefaultConfig)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
