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
			analyzer.ConstantKindFlag:       cfg.ConstantKind,
			analyzer.CryptoHashFlag:         cfg.CryptoHash,
			analyzer.HTTPMethodFlag:         cfg.HTTPMethod,
			analyzer.HTTPStatusCodeFlag:     cfg.HTTPStatusCode,
			analyzer.OSDevNullFlag:          cfg.OSDevNullFlag,
			analyzer.RPCDefaultPathFlag:     cfg.DefaultRPCPathFlag,
			analyzer.SQLIsolationLevelFlag:  cfg.SQLIsolationLevelFlag,
			analyzer.TimeLayoutFlag:         cfg.TimeLayout,
			analyzer.TimeMonthFlag:          cfg.TimeMonth,
			analyzer.TimeWeekdayFlag:        cfg.TimeWeekday,
			analyzer.TLSSignatureSchemeFlag: cfg.TLSSignatureSchemeFlag,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
