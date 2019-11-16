package golinters

import (
	magic_numbers "github.com/tommy-muehle/go-mnd"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGomnd() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		magic_numbers.Analyzer,
	}

	return goanalysis.NewLinter(
		"gomnd",
		"checks whether magic number is used",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
