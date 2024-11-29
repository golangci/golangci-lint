package gosmopolitan

import (
	"strings"

	"github.com/xen0n/gosmopolitan"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.GosmopolitanSettings) *goanalysis.Linter {
	a := gosmopolitan.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"allowtimelocal":  settings.AllowTimeLocal,
			"escapehatches":   strings.Join(settings.EscapeHatches, ","),
			"lookattests":     !settings.IgnoreTests,
			"watchforscripts": strings.Join(settings.WatchForScripts, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
