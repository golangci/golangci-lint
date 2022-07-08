package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/gostaticanalysis/forcetypeassert"
	"golang.org/x/tools/go/analysis"
)

func NewForceTypeAssert() *goanalysis.Linter {
	a := forcetypeassert.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"finds forced type assertions",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
