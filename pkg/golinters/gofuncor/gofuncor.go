package gofuncor

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/manuelarte/gofuncor/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New() *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
