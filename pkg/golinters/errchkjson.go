package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/breml/errchkjson"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewErrChkJSONFuncName(cfg *config.ErrChkJSONSettings) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil {
		if cfg.OmitSafe {
			cfgMap[a.Name] = map[string]interface{}{
				"omit-safe": "true",
			}
		}
	}

	return goanalysis.NewLinter(
		"errchkjson",
		"Checks types passed to the json encoding functions. "+
			"Reports unsupported types and reports occations, where the check for the returned error can be omitted.",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
