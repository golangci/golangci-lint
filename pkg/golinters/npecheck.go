package golinters

import (
	"golang.org/x/tools/go/analysis"

	check "github.com/chenfeining/go-npecheck"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNilPointerReferenceCheck() *goanalysis.Linter {
	analyzer := check.Analyzer
	return goanalysis.NewLinter(
		analyzer.Name,
		check.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)
}
