package tagliatelle

import (
	"github.com/ldez/tagliatelle"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.TagliatelleSettings) *goanalysis.Linter {
	cfg := tagliatelle.Config{
		Base: tagliatelle.Base{
			Rules: map[string]string{
				"json":   "camel",
				"yaml":   "camel",
				"header": "header",
			},
		},
	}

	if settings != nil {
		for k, v := range settings.Case.Rules {
			cfg.Rules[k] = v
		}

		cfg.ExtendedRules = toExtendedRules(settings.Case.ExtendedRules)
		cfg.UseFieldName = settings.Case.UseFieldName
		cfg.IgnoredFields = settings.Case.IgnoredFields

		for _, override := range settings.Case.Overrides {
			cfg.Overrides = append(cfg.Overrides, tagliatelle.Overrides{
				Base: tagliatelle.Base{
					Rules:         override.Rules,
					ExtendedRules: toExtendedRules(override.ExtendedRules),
					UseFieldName:  override.UseFieldName,
					IgnoredFields: override.IgnoredFields,
					Ignore:        override.Ignore,
				},
				Package: override.Package,
			})
		}
	}

	a := tagliatelle.New(cfg)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func toExtendedRules(src map[string]config.TagliatelleExtendedRule) map[string]tagliatelle.ExtendedRule {
	result := make(map[string]tagliatelle.ExtendedRule, len(src))

	for k, v := range src {
		result[k] = tagliatelle.ExtendedRule{
			Case:                v.Case,
			ExtraInitialisms:    v.ExtraInitialisms,
			InitialismOverrides: v.InitialismOverrides,
		}
	}

	return result
}
