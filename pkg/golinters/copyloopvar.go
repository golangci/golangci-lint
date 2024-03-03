package golinters

import (
	"github.com/karamaru-alpha/copyloopvar"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewCopyLoopVar(settings *config.CopyLoopVarSettings) *goanalysis.Linter {
	analyzer := copyloopvar.NewAnalyzer()
	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			analyzer.Name: {
				"ignore-alias": settings.IgnoreAlias,
			},
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
