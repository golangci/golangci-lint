package golinters

import (
	"github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewUseStdlibVars(cfg *config.UseStdlibVarsSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]interface{})
	if cfg != nil {
		cfgMap[a.Name] = map[string]interface{}{
			analyzer.HTTPMethodFlag:     cfg.HTTPMethod,
			analyzer.HTTPStatusCodeFlag: cfg.HTTPStatusCode,
			analyzer.TimeWeekdayFlag:    cfg.TimeWeekday,
			analyzer.TimeMonthFlag:      cfg.TimeMonth,
			analyzer.TimeLayoutFlag:     cfg.TimeLayout,
			analyzer.CryptoHashFlag:     cfg.CryptoHash,
			analyzer.DefaultRPCPathFlag: cfg.DefaultRPCPathFlag,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
