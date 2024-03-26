package golinters

import (
	"strings"

	"github.com/Abirdcfly/dupword"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func NewDupWord(setting *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	if setting != nil {
		cfgMap[a.Name] = map[string]any{
			"keyword": strings.Join(setting.Keywords, ","),
			"ignore":  strings.Join(setting.Ignore, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"checks for duplicate words in the source code",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
