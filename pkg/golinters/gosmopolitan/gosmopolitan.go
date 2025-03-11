package gosmopolitan

import (
	"strings"

	"github.com/xen0n/gosmopolitan"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GosmopolitanSettings) *goanalysis.Linter {
	a := gosmopolitan.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"allowtimelocal":  settings.AllowTimeLocal,
			"escapehatches":   strings.Join(settings.EscapeHatches, ","),
			"watchforscripts": strings.Join(settings.WatchForScripts, ","),

			// Should be managed with `linters.exclusions.rules`.
			"lookattests": true,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
