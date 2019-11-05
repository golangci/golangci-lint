package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"

	magic_numbers "github.com/tommy-muehle/go-mnd"
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
