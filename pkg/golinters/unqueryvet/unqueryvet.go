package unqueryvet

import (
	"github.com/MirrexOne/unqueryvet"
	pkgconfig "github.com/MirrexOne/unqueryvet/pkg/config"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.UnqueryvetSettings) *goanalysis.Linter {
	cfg := pkgconfig.DefaultSettings()

	if settings != nil {
		cfg.CheckSQLBuilders = settings.CheckSQLBuilders
		if len(settings.AllowedPatterns) > 0 {
			cfg.AllowedPatterns = settings.AllowedPatterns
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(unqueryvet.NewWithConfig(&cfg)).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
