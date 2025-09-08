package gounqvet

import (
	"github.com/MirrexOne/gounqvet"
	pkgconfig "github.com/MirrexOne/gounqvet/pkg/config"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GounqvetSettings) *goanalysis.Linter {
	cfg := pkgconfig.DefaultSettings()

	if settings != nil {
		cfg.CheckSQLBuilders = settings.CheckSQLBuilders
		if len(settings.AllowedPatterns) > 0 {
			cfg.AllowedPatterns = settings.AllowedPatterns
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(gounqvet.NewWithConfig(&cfg)).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
