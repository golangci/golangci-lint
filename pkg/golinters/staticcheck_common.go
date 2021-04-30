package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

var debugf = logutils.Debug("megacheck")

func setupStaticCheckAnalyzers(m map[string]*analysis.Analyzer, settings *config.StaticCheckSettings) []*analysis.Analyzer {
	var ret []*analysis.Analyzer
	for _, v := range m {
		setAnalyzerGoVersion(v, settings)
		ret = append(ret, v)
	}
	return ret
}

func setAnalyzerGoVersion(a *analysis.Analyzer, settings *config.StaticCheckSettings) {
	// TODO: uses "1.13" for backward compatibility, but in the future (v2) must be set by using build.Default.ReleaseTags like staticcheck.
	goVersion := "1.13"
	if settings != nil && settings.GoVersion != "" {
		goVersion = settings.GoVersion
	}

	if v := a.Flags.Lookup("go"); v != nil {
		if err := v.Value.Set(goVersion); err != nil {
			debugf("Failed to set go version: %s", err)
		}
	}
}
