package golinters

import (
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewRowsErrCheck(settings *config.RowsErrCheckSettings) *goanalysis.Linter {
	var pkgs []string
	if settings != nil {
		pkgs = settings.Packages
	}

	analyzer := rowserr.NewAnalyzer(pkgs...)

	return goanalysis.NewLinter(
		"rowserrcheck",
		"checks whether Err of rows is checked successfully",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
