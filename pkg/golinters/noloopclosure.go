package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/fatanugraha/noloopclosure"
)

func NewNoLoopClosure() *goanalysis.Linter {
	analyzer := noloopclosure.Analyzer

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
