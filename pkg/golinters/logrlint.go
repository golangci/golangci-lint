package golinters

import (
	"strconv"
	"strings"

	"github.com/timonwong/logrlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewLogrLint(settings *config.LogrLintSettings) *goanalysis.Linter {
	analyzer := logrlint.NewAnalyzer()
	cfg := map[string]map[string]interface{}{}
	if settings != nil {
		linterCfg := map[string]interface{}{
			"disableall": strconv.FormatBool(settings.DisableAll),
			"enable":     strings.Join(settings.Enable, ","),
			"disable":    strings.Join(settings.Disable, ","),
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
