package exhaustive

import (
	"strconv"

	"github.com/nishanths/exhaustive/passes/enumerated"
	"github.com/nishanths/exhaustive/passes/exhaustive"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.ExhaustiveSettings) *goanalysis.Linter {
	exhAnalyzer := exhaustive.Analyzer

	var exhCfg, enuCfg map[string]any

	if settings != nil {
		if settings.ExplicitExhaustiveSwitch {
			internal.LinterLogger.Warnf("%s: the option 'explicit-exhaustive-switch' has been replaced by 'explicit-exhaustive'.", exhAnalyzer.Name)
		}

		if settings.ExplicitExhaustiveMap {
			internal.LinterLogger.Warnf("%s: the option 'explicit-exhaustive-map' has been replaced by 'explicit-exhaustive'.", exhAnalyzer.Name)
		}

		exhCfg = map[string]any{
			"d":          settings.DefaultSignifiesExhaustive,
			"e":          settings.ExplicitExhaustive || settings.ExplicitExhaustiveSwitch || settings.ExplicitExhaustiveMap,
			"defrequire": settings.DefaultCaseRequired,
			// Should be managed with `linters.exclusions.generated`.
			"g": true,
		}

		if len(settings.Check) > 0 {
			var checks []string

			for _, s := range settings.Check {
				if s == "map" {
					internal.LinterLogger.Warnf("%s: the value 'map' has been replaced by 'mapliteral'.", exhAnalyzer.Name)

					checks = append(checks, "mapliteral")
				} else {
					checks = append(checks, s)
				}
			}

			exhCfg["check"] = checks
		}

		repeatedFlags(exhCfg, "typeonly", settings.OnlyTypes)

		// An empty string is converted to a regular expression that exclude everything.
		if settings.IgnoreEnumMembers != "" {
			internal.LinterLogger.Warnf("%s: the option 'ignore-enum-members' has been replaced by 'ignore-constants'.", exhAnalyzer.Name)

			exhCfg["constignore"] = settings.IgnoreEnumMembers
		}

		repeatedFlags(exhCfg, "constignore", settings.IgnoreConstants)

		// An empty string is converted to a regular expression that exclude everything.
		if settings.IgnoreEnumTypes != "" {
			internal.LinterLogger.Warnf("%s: the option 'ignore-enum-types' has been replaced by 'ignore-types'.", exhAnalyzer.Name)

			exhCfg["typeignore"] = settings.IgnoreEnumTypes
		}

		repeatedFlags(exhCfg, "typeignore", settings.IgnoreTypes)

		enuCfg = map[string]any{
			"p": settings.PackageScopeOnly,
		}
	}

	return goanalysis.
		NewLinter(
			exhAnalyzer.Name,
			exhAnalyzer.Doc,
			[]*analysis.Analyzer{
				exhAnalyzer,
				enumerated.Analyzer,
			},
			map[string]map[string]any{
				enumerated.Analyzer.Name: enuCfg,
				exhAnalyzer.Name:         exhCfg,
			},
		).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func repeatedFlags(cfg map[string]any, prefix string, values []string) {
	for i, part := range values {
		cfg[prefix+"-"+strconv.Itoa(i)+"$"] = part
	}
}
