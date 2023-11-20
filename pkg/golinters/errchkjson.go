package golinters

import (
	"github.com/breml/errchkjson"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewErrChkJSONFuncName(cfg *config.ErrChkJSONSettings) *goanalysis.Linter {
	a := errchkjson.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	cfgMap[a.Name] = map[string]any{
		"omit-safe": true,
	}
	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			"omit-safe":          !cfg.CheckErrorFreeEncoding,
			"report-no-exported": cfg.ReportNoExported,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
