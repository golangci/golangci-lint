package golinters

import (
	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGocheckcompilerdirectives() *goanalysis.Linter {
	gocheckcompilerdirectives := checkcompilerdirectives.Analyzer()

	return goanalysis.NewLinter(
		gocheckcompilerdirectives.Name,
		gocheckcompilerdirectives.Doc,
		[]*analysis.Analyzer{gocheckcompilerdirectives},
		linterConfig,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
