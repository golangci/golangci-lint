package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/katsadim/nostdglobals"
	"golang.org/x/tools/go/analysis"
)

func NewNoStdGlobals() *goanalysis.Linter {
	a := nostdglobals.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
