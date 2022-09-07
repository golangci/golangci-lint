package golinters

import (
	"github.com/dorfire/go-analyzers/src/onlyany"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewOnlyAny() *goanalysis.Linter {
	return goanalysis.NewLinter(
		onlyany.Analyzer.Name,
		onlyany.Analyzer.Doc,
		[]*analysis.Analyzer{onlyany.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
