package canonicalheader

import (
	"strings"

	"github.com/golangci/canonicalheader"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.CanonicalHeaderSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"useDefaultExclusion": settings.UseDefaultExclusions,
		}

		if len(settings.Exclusions) > 0 {
			cfg["exclusions"] = strings.Join(settings.Exclusions, ",")
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(canonicalheader.New()).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
