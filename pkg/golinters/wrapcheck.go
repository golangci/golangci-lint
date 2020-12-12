package golinters

import (
	"github.com/tomarrell/wrapcheck/wrapcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
)

const wrapcheckName = "wrapcheck"

func NewWrapcheck() *goanalysis.Linter {
	return goanalysis.NewLinter(
		wrapcheckName,
		wrapcheck.Analyzer.Doc,
		[]*analysis.Analyzer{wrapcheck.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
