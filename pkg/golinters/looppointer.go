package golinters

import (
	"github.com/kyoh86/looppointer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

// NewLoopPointer returns a new instance of Linter based on
// github.com/kyoh86/looppointer.
func NewLoopPointer() *goanalysis.Linter {
	analyzer := looppointer.Analyzer

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
