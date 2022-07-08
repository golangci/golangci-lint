package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/kyoh86/exportloopref"
	"golang.org/x/tools/go/analysis"
)

func NewExportLoopRef() *goanalysis.Linter {
	a := exportloopref.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
