package golinters

import (
	"strings"

	"github.com/timonwong/todolint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewTODOLint(settings *config.TODOLintSettings) *goanalysis.Linter {
	analyzer := todolint.NewAnalyzer()

	cfgMap := map[string]map[string]interface{}{}
	if settings != nil {
		keywords := strings.Join(settings.Keywords, ",")
		cfg := map[string]interface{}{}
		if keywords != "" {
			cfg["keywords"] = keywords
		}

		cfgMap[analyzer.Name] = cfg
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
