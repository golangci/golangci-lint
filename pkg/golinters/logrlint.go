package golinters

import (
	"github.com/timonwong/logrlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const LogrLintName = "logrlint"

func NewLogrLint() *goanalysis.Linter {
	return goanalysis.NewLinter(
		LogrLintName,
		logrlint.Doc,
		[]*analysis.Analyzer{logrlint.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
