package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/sivchari/nosnakecase"
	"golang.org/x/tools/go/analysis"
)

func NewNoSnakeCase() *goanalysis.Linter {
	a := nosnakecase.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
