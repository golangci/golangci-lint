package golinters

import (
	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoCheckCompilerDirectives() *goanalysis.Linter {
	a := checkcompilerdirectives.Analyzer()

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
