package pkgname

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/uudashr/pkgname"
)

func New(settings *config.PkgnameSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"include-import-alias": settings.ImportAlias,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(pkgname.Analyzer).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
