package tfproviderlint

import (
	"slices"
	"strings"

	"github.com/bflad/tfproviderlint/passes"
	"github.com/bflad/tfproviderlint/xpasses"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.TFProviderLintSettings) *goanalysis.Linter {
	var conf map[string]map[string]any

	if settings != nil {
		conf = make(map[string]map[string]any)

		for k, v := range settings.Settings {
			// The settings related to AT001 and AT012 are ignored because they must be handled globally with exclusions.
			if k == "AT001" || k == "AT012" {
				continue
			}

			// The analyzer must be in uppercase, but Viper converts them to lowercase.
			conf[strings.ToUpper(k)] = v
		}
	}

	return goanalysis.NewLinter(
		"tfproviderlint",
		"Linter for Terraform Provider code.",
		analyzersFromSettings(settings),
		conf,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromSettings(settings *config.TFProviderLintSettings) []*analysis.Analyzer {
	if settings == nil {
		return passes.AllChecks
	}

	switch settings.Default {
	case "standard", "":
		return modeAll(settings, false)

	case "extra":
		return modeAll(settings, true)

	case "none":
		return modeNone(settings)

	default:
		return nil
	}
}

func modeAll(settings *config.TFProviderLintSettings, extra bool) []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer

	if len(settings.Disable) == 0 {
		if extra {
			return slices.Concat(passes.AllChecks, xpasses.AllChecks)
		}

		return passes.AllChecks
	}

	for _, check := range passes.AllChecks {
		if !slices.Contains(settings.Disable, check.Name) {
			analyzers = append(analyzers, check)
		}
	}

	if extra {
		for _, check := range xpasses.AllChecks {
			if !slices.Contains(settings.Enable, check.Name) {
				analyzers = append(analyzers, check)
			}
		}
	}

	return analyzers
}

func modeNone(settings *config.TFProviderLintSettings) []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer

	all := slices.Concat(passes.AllChecks, xpasses.AllChecks)

	for _, check := range all {
		if slices.Contains(settings.Enable, check.Name) {
			analyzers = append(analyzers, check)
		}
	}

	return analyzers
}
