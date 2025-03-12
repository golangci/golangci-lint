package forcetypeassert

import (
	"github.com/gostaticanalysis/forcetypeassert"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := forcetypeassert.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"finds forced type assertions",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
