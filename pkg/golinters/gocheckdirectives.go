package golinters

import (
	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGocheckcompilerdirectives() *goanalysis.Linter {
	gocheckcompilerdirectives := checkcompilerdirectives.Analyzer()

	// gocheckcompilerdirectives has no config.
	linterConfig := map[string]map[string]interface{}{
		gocheckcompilerdirectives.Name: {},
	}

	return goanalysis.NewLinter(
		gocheckcompilerdirectives.Name,
		gocheckcompilerdirectives.Doc,
		[]*analysis.Analyzer{gocheckcompilerdirectives},
		linterConfig,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
