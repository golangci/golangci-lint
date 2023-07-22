package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/nilpointer"
)

func NewNilPointerReferenceCheck() *goanalysis.Linter {
	a := nilpointer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		nilpointer.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)
}
