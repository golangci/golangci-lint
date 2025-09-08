package gounqvet

import (
	"github.com/MirrexOne/gounqvet"
	pkgconfig "github.com/MirrexOne/gounqvet/pkg/config"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GounqvetSettings) *goanalysis.Linter {
	var cfg *pkgconfig.GounqvetSettings

	if settings != nil {
		cfg = &pkgconfig.GounqvetSettings{
			CheckSQLBuilders: settings.CheckSQLBuilders,
			AllowedPatterns:  settings.AllowedPatterns,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(gounqvet.NewWithConfig(cfg)).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
