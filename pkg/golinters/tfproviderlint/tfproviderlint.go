package tfproviderlint

import (
	"slices"
	"strconv"

	"github.com/bflad/tfproviderlint/passes"
	"github.com/bflad/tfproviderlint/passes/AT001"
	"github.com/bflad/tfproviderlint/passes/AT012"
	"github.com/bflad/tfproviderlint/passes/R006"
	"github.com/bflad/tfproviderlint/passes/R019"
	"github.com/bflad/tfproviderlint/xpasses"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func New(settings *config.TfproviderlintSettings) *goanalysis.Linter {
	analyzers := analyzersFromSettings(settings)

	return goanalysis.
		NewLinter("tfproviderlint", "Static analysis for Terraform Provider code (resource, schema, acceptance tests)", analyzers, nil).
		WithContextSetter(func(lintCtx *linter.Context) {
			if settings == nil {
				return
			}

			// AT001 settings
			if settings.AT001.IgnoredFilenamePrefixes != "" {
				if err := AT001.Analyzer.Flags.Set("ignored-filename-prefixes", settings.AT001.IgnoredFilenamePrefixes); err != nil {
					lintCtx.Log.Errorf("tfproviderlint: failed to set AT001.ignored-filename-prefixes: %v", err)
				}
			}
			if settings.AT001.IgnoredFilenameSuffixes != "" {
				if err := AT001.Analyzer.Flags.Set("ignored-filename-suffixes", settings.AT001.IgnoredFilenameSuffixes); err != nil {
					lintCtx.Log.Errorf("tfproviderlint: failed to set AT001.ignored-filename-suffixes: %v", err)
				}
			}

			// AT012 settings
			if settings.AT012.IgnoredFilenames != "" {
				if err := AT012.Analyzer.Flags.Set("ignored-filenames", settings.AT012.IgnoredFilenames); err != nil {
					lintCtx.Log.Errorf("tfproviderlint: failed to set AT012.ignored-filenames: %v", err)
				}
			}

			// R006 settings
			if settings.R006.PackageAliases != "" {
				if err := R006.Analyzer.Flags.Set("package-aliases", settings.R006.PackageAliases); err != nil {
					lintCtx.Log.Errorf("tfproviderlint: failed to set R006.package-aliases: %v", err)
				}
			}

			// R019 settings
			if settings.R019.Threshold != 0 {
				if err := R019.Analyzer.Flags.Set("threshold", strconv.Itoa(settings.R019.Threshold)); err != nil {
					lintCtx.Log.Errorf("tfproviderlint: failed to set R019.threshold: %v", err)
				}
			}
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromSettings(settings *config.TfproviderlintSettings) []*analysis.Analyzer {
	// Build list of analyzers based on settings
	var allAnalyzers []*analysis.Analyzer

	if settings == nil {
		// Default: all standard and extra passes
		return slices.Concat(passes.AllChecks, xpasses.AllChecks)
	}

	// Add standard passes if enabled (default: true)
	if settings.IsAllEnabled() {
		allAnalyzers = append(allAnalyzers, passes.AllChecks...)
	}

	// Add extra passes if enabled (default: true)
	if settings.IsExtraEnabled() {
		allAnalyzers = append(allAnalyzers, xpasses.AllChecks...)
	}

	// Filter based on enable/disable lists
	var enabledAnalyzers []*analysis.Analyzer
	for _, a := range allAnalyzers {
		if isAnalyzerEnabled(a.Name, settings) {
			enabledAnalyzers = append(enabledAnalyzers, a)
		}
	}

	return enabledAnalyzers
}

func isAnalyzerEnabled(name string, cfg *config.TfproviderlintSettings) bool {
	switch {
	case cfg.DisableAll:
		return slices.Contains(cfg.Enable, name)

	case slices.Contains(cfg.Disable, name):
		return false

	case slices.Contains(cfg.Enable, name):
		return true

	default:
		return true
	}
}
