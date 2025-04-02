package nakedret

import (
	"github.com/alexkohler/nakedret/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.NakedretSettings) *goanalysis.Linter {
	cfg := &nakedret.NakedReturnRunner{}

	if settings != nil {
		// SkipTestFiles is intentionally ignored => should be managed with `linters.exclusions.rules`.
		cfg.MaxLength = settings.MaxFuncLines
	}

	a := nakedret.NakedReturnAnalyzer(cfg)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
