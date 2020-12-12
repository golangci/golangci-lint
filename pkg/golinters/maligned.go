package golinters

import (
	"github.com/mdempsky/maligned/passes/maligned"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewMaligned() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"maligned",
		maligned.Analyzer.Doc,
		[]*analysis.Analyzer{maligned.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
