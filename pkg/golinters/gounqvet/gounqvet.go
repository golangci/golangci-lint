package gounqvet

import (
	"golang.org/x/tools/go/analysis"

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

	analyzer := gounqvet.NewWithConfig(cfg)
	
	return goanalysis.NewLinter(
		"gounqvet",
		"Detects SELECT * usage in SQL queries and SQL builders, encouraging explicit column selection",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
