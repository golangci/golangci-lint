package golinters

import (
	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoctx() *goanalysis.Linter {
	a := noctx.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		"Detects test helpers which is not start with t.Helper() method",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
