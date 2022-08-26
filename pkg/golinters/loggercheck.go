package golinters

import (
	"strings"

	"github.com/timonwong/loggercheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewLoggerCheck(settings *config.LoggerCheckSettings) *goanalysis.Linter {
	analyzer := loggercheck.NewAnalyzer()
	cfg := map[string]map[string]interface{}{}
	if settings != nil {
		var disabled []string
		if !settings.Logr {
			disabled = append(disabled, "logr")
		}
		if !settings.Klog {
			disabled = append(disabled, "klog")
		}
		if !settings.Logr {
			disabled = append(disabled, "zap")
		}
		linterCfg := map[string]interface{}{
			"disable": strings.Join(disabled, ","),
		}
		cfg[analyzer.Name] = linterCfg
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
