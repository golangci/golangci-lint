package golinters

import (
	"github.com/timonwong/loggercheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewLoggerCheck(settings *config.LoggerCheckSettings) *goanalysis.Linter {
	var (
		disable []string
		rules   []string
	)

	if settings != nil {
		if !settings.Logr {
			disable = append(disable, "logr")
		}
		if !settings.Klog {
			disable = append(disable, "klog")
		}
		if !settings.Zap {
			disable = append(disable, "zap")
		}

		rules = settings.Rules
	}

	analyzer := loggercheck.NewAnalyzer(loggercheck.WithDisable(disable),
		loggercheck.WithRules(rules))
	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
