package docenv

import (
	"github.com/g4s8/envdoc/linter"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(cfg *config.DocenvSettings) *goanalysis.Linter {
	opts := []linter.Option{linter.WithNoComments()}
	if cfg != nil {
		opts = append(opts, linter.WithEnvName(cfg.TagName))
	}

	a := linter.NewAnlyzer(false, opts...)
	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
