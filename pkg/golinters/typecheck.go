package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewTypecheck() *goanalysis.Linter {
	const linterName = "typecheck"

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return nil, nil
		},
	}

	// Note: typecheck doesn't require the LoadModeWholeProgram
	// but it's a hack to force this linter to be the first linter in all the cases.
	linter := goanalysis.NewLinter(
		linterName,
		"Like the front-end of a Go compiler, parses and type-checks Go code",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)

	linter.SetTypecheckMode()

	return linter
}
