package usestdlibvars

import (
	"github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.UseStdlibVarsSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			analyzer.ConstantKindFlag:       settings.ConstantKind,
			analyzer.CryptoHashFlag:         settings.CryptoHash,
			analyzer.HTTPMethodFlag:         settings.HTTPMethod,
			analyzer.HTTPStatusCodeFlag:     settings.HTTPStatusCode,
			analyzer.OSDevNullFlag:          settings.OSDevNull != nil && *settings.OSDevNull,
			analyzer.RPCDefaultPathFlag:     settings.DefaultRPCPath,
			analyzer.SQLIsolationLevelFlag:  settings.SQLIsolationLevel,
			analyzer.SyslogPriorityFlag:     settings.SyslogPriority != nil && *settings.SyslogPriority,
			analyzer.TimeLayoutFlag:         settings.TimeLayout,
			analyzer.TimeMonthFlag:          settings.TimeMonth,
			analyzer.TimeWeekdayFlag:        settings.TimeWeekday,
			analyzer.TLSSignatureSchemeFlag: settings.TLSSignatureScheme,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
