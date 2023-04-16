package golinters

import (
	"github.com/ykadowak/zerologlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewZerologLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"zerologlint",
		"Detects the wrong usage of `zerolog` that a user forgets to dispatch with `Send` or `Msg`.",
		[]*analysis.Analyzer{zerologlint.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
