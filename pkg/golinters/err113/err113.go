package err113

import (
	"github.com/Djarvur/go-err113"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	a := err113.NewAnalyzer()

	return goanalysis.NewLinter(
		a.Name,
		"Go linter to check the errors handling expressions",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
