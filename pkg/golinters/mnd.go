package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/tommy-muehle/go-mnd"
)

func NewMnd() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		magic_numbers.Analyzer,
	}

	return goanalysis.NewLinter(
		"magicnumber",
		"checks whether magic number detector is used",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
