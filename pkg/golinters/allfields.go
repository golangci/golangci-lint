package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/subtle-byte/allfields/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewAllfields() *goanalysis.Linter {
	a := analyzer.Analyzer
	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
