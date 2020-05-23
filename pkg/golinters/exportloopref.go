package golinters

import (
	"github.com/kyoh86/exportloopref"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExportLoopRef() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		exportloopref.Analyzer,
	}

	return goanalysis.NewLinter(
		"exportloopref",
		"An analyzer that finds exporting pointers for loop variables.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
