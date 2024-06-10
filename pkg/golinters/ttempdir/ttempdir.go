package ttempdir

import (
	"github.com/peczenyj/ttempdir/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TtempdirSettings) *goanalysis.Linter {
	a := analyzer.New()

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				analyzer.FlagAllName:               settings.All,
				analyzer.FlagMaxRecursionLevelName: settings.MaxRecursionLevel,
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"Detects the use of os.MkdirTemp, ioutil.TempDir or os.TempDir instead of t.TempDir",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
