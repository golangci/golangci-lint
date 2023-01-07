package golinters

import (
	"github.com/leighmcculloch/gocheckdirectives/checkdirectives"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGocheckdirectives() *goanalysis.Linter {
	gocheckdirectives := checkdirectives.Analyzer()

	// gocheckdirectives has no config.
	linterConfig := map[string]map[string]interface{}{
		gocheckdirectives.Name: {},
	}

	return goanalysis.NewLinter(
		gocheckdirectives.Name,
		gocheckdirectives.Doc,
		[]*analysis.Analyzer{gocheckdirectives},
		linterConfig,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
