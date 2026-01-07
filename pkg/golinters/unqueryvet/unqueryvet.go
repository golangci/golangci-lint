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
		// IgnoredFiles, and Severity are explicitly ignored.
		cfg.CheckSQLBuilders = settings.CheckSQLBuilders
		cfg.CheckAliasedWildcard = settings.CheckAliasedWildcard
		cfg.CheckStringConcat = settings.CheckStringConcat
		cfg.CheckFormatStrings = settings.CheckFormatStrings
		cfg.CheckStringBuilder = settings.CheckStringBuilder
		cfg.CheckSubqueries = settings.CheckSubqueries
		cfg.IgnoredFunctions = settings.IgnoredFunctions

		if len(settings.AllowedPatterns) > 0 {
			cfg.AllowedPatterns = settings.AllowedPatterns
		}

		cfg.SQLBuilders = pkgconfig.SQLBuildersConfig{
			Squirrel:  settings.SQLBuilders.Squirrel,
			GORM:      settings.SQLBuilders.GORM,
			SQLx:      settings.SQLBuilders.SQLx,
			Ent:       settings.SQLBuilders.Ent,
			PGX:       settings.SQLBuilders.PGX,
			Bun:       settings.SQLBuilders.Bun,
			SQLBoiler: settings.SQLBuilders.SQLBoiler,
			Jet:       settings.SQLBuilders.Jet,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(unqueryvet.NewWithConfig(&cfg)).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
