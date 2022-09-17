package golinters

import (
	"github.com/maratori/pairedbrackets/pkg/pairedbrackets"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewPairedbrackets(cfg *config.PairedBracketsSettings) *goanalysis.Linter {
	var a = pairedbrackets.NewAnalyzer()

	var settings map[string]map[string]interface{}
	if cfg != nil {
		settings = map[string]map[string]interface{}{
			a.Name: {
				pairedbrackets.IgnoreFuncCallsFlagName: cfg.IgnoreFuncCalls,
			},
		}
	}

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, settings).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
