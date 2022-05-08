package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/dgunay/ifacecapture/ifacecapture"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewInterfaceCaptureCheck(settings *config.IfaceCaptureSettings) *goanalysis.Linter {
	linterCfg := map[string]map[string]interface{}{}
	if settings != nil {
		linterCfg["grouper"] = map[string]interface{}{
			"loglvl":            settings.LogLevel,
			"ignore-interfaces": settings.IgnoreInterfaces,
			"allow-interfaces":  settings.AllowInterfaces,
		}
	}

	ifacecaptureAnalyzer := ifacecapture.Analyzer
	return goanalysis.NewLinter(
		ifacecaptureAnalyzer.Name,
		ifacecaptureAnalyzer.Doc,
		[]*analysis.Analyzer{ifacecaptureAnalyzer},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
