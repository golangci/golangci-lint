package gosmopolitan

import (
	"strings"

	"github.com/xen0n/gosmopolitan"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(s *config.GosmopolitanSettings) *goanalysis.Linter {
	a := gosmopolitan.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	if s != nil {
		cfgMap[a.Name] = map[string]any{
			"allowtimelocal":  s.AllowTimeLocal,
			"escapehatches":   strings.Join(s.EscapeHatches, ","),
			"lookattests":     !s.IgnoreTests,
			"watchforscripts": strings.Join(s.WatchForScripts, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
