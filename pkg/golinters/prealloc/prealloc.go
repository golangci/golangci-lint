package prealloc

import (
	"fmt"

	"github.com/alexkohler/prealloc/pkg"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

const linterName = "prealloc"

func New(settings *config.PreallocSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			runPreAlloc(pass, settings)

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Finds slice declarations that could potentially be pre-allocated",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPreAlloc(pass *analysis.Pass, settings *config.PreallocSettings) {
	hints := pkg.Check(pass.Files, settings.Simple, settings.RangeLoops, settings.ForLoops)

	for _, hint := range hints {
		pass.Report(analysis.Diagnostic{
			Pos:     hint.Pos,
			Message: fmt.Sprintf("Consider pre-allocating %s", internal.FormatCode(hint.DeclaredSliceName, nil)),
		})
	}
}
