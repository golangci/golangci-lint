package docnametypo

import (
	"github.com/cce/docnametypo/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.DocnameTypoSettings) *goanalysis.Linter {
	a := analyzer.Analyzer

	cfg := map[string]map[string]any{}
	if settings != nil {
		linterCfg := map[string]any{
			"maxdist":                   settings.MaxDist,
			"include-unexported":        settings.IncludeUnexported,
			"include-exported":          settings.IncludeExported,
			"include-types":             settings.IncludeTypes,
			"include-generated":         settings.IncludeGenerated,
			"include-interface-methods": settings.IncludeInterfaceMethods,
		}

		if settings.AllowedLeadingWords != "" {
			linterCfg["allowed-leading-words"] = settings.AllowedLeadingWords
		}

		if settings.AllowedPrefixes != "" {
			linterCfg["allowed-prefixes"] = settings.AllowedPrefixes
		}

		cfg[a.Name] = linterCfg
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*goanalysis.Analyzer{{Analyzer: a}},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
