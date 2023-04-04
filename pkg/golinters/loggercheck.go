package golinters

import (
	"github.com/timonwong/loggercheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewLoggerCheck(settings *config.LoggerCheckSettings) *goanalysis.Linter {
	var opts []loggercheck.Option

	if settings != nil {
		var disable []string
		if !settings.Kitlog {
			disable = append(disable, "kitlog")
		}
		if !settings.Klog {
			disable = append(disable, "klog")
		}
		if !settings.Logr {
			disable = append(disable, "logr")
		}
		if !settings.Zap {
			disable = append(disable, "zap")
		}

		opts = []loggercheck.Option{
			loggercheck.WithDisable(disable),
			loggercheck.WithRequireStringKey(settings.RequireStringKey),
			loggercheck.WithRules(settings.Rules),
			loggercheck.WithNoPrintfLike(settings.NoPrintfLike),
		}
	}

	analyzer := loggercheck.NewAnalyzer(opts...)
	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
