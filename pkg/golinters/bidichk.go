package golinters

import (
	"github.com/breml/bidichk/pkg/bidichk"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewBiDiChkFuncName() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"bidichk",
		"Checks for dangerous unicode character sequences",
		[]*analysis.Analyzer{bidichk.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
