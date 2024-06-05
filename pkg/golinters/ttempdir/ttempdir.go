package ttempdir

import (
	"github.com/peczenyj/ttempdir/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TtempdirSettings) *goanalysis.Linter {
	ttempdirAnalyzer := analyzer.New()

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			ttempdirAnalyzer.Name: {
				analyzer.FlagAllName:               settings.All,
				analyzer.FlagMaxRecursionLevelName: settings.MaxRecursionLevel,
			},
		}
	}

	return goanalysis.NewLinter(
		ttempdirAnalyzer.Name,
		ttempdirAnalyzer.Doc,
		[]*analysis.Analyzer{ttempdirAnalyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
