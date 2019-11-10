package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"

	magic_numbers "github.com/tommy-muehle/go-mnd"
)

func NewGomnd() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		magic_numbers.Analyzer,
	}

	return goanalysis.NewLinter(
		"gomnd",
		"checks whether magic number detector is used",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
