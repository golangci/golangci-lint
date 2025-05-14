package prealloc

import (
	"fmt"

	"github.com/alexkohler/prealloc/pkg"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.PreallocSettings) *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: "prealloc",
			Doc:  "Find slice declarations that could potentially be pre-allocated",
			Run: func(pass *analysis.Pass) (any, error) {
				runPreAlloc(pass, settings)

				return nil, nil
			},
		}).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func runPreAlloc(pass *analysis.Pass, settings *config.PreallocSettings) {
	hints := pkg.Check(pass.Files, settings.Simple, settings.RangeLoops, settings.ForLoops)

	for _, hint := range hints {
		pass.Report(analysis.Diagnostic{
			Pos:     hint.Pos,
			Message: fmt.Sprintf("Consider pre-allocating %s", internal.FormatCode(hint.DeclaredSliceName)),
		})
	}
}
