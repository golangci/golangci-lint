package golinters

import (
	"strings"

	"github.com/Abirdcfly/dupword"
	"github.com/golangci/golangci-lint/pkg/config"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewDupWord(setting *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()
	cfgMap := map[string]map[string]interface{}{}
	if setting != nil {
		cfgMap[a.Name] = map[string]interface{}{
			"keyword": strings.Join(setting.KeyWord, ","),
		}
	}
	return goanalysis.NewLinter(
		"dupword",
		"checks for duplicate words in the source code (usually miswritten)",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
