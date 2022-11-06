package golinters

import (
	"github.com/stevenh/go-rowserr/pkg/rowserr"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewRowsErr(settings *config.RowsErrSettings) *goanalysis.Linter {
	a := rowserr.NewAnalyzer(settings.Packages...)
	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
