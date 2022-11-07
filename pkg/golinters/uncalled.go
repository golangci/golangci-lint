package golinters

import (
	"github.com/stevenh/go-uncalled/pkg/uncalled"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewUncalled(settings *config.UncalledSettings) *goanalysis.Linter {
	a := uncalled.NewAnalyzer(settings)
	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
