package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"go.uber.org/nilaway"
	"golang.org/x/tools/go/analysis"
)

func NewNilAway() *goanalysis.Linter {
	a := nilaway.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
