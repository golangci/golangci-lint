package iface

import (
	"slices"
	"strings"

	"github.com/uudashr/iface/identical"
	"github.com/uudashr/iface/opaque"
	"github.com/uudashr/iface/unused"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.IfaceSettings) *goanalysis.Linter {
	var conf map[string]map[string]any
	if settings != nil {
		conf = settings.Settings
	}

	return goanalysis.NewLinter(
		"iface",
		"Detect the incorrect use of interfaces, helping developers avoid interface pollution.",
		analyzersFromSettings(settings),
		conf,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromSettings(settings *config.IfaceSettings) []*analysis.Analyzer {
	allAnalyzers := map[string]*analysis.Analyzer{
		"unused":    unused.Analyzer,
		"identical": identical.Analyzer,
		"opaque":    opaque.Analyzer,
	}

	if settings == nil || len(settings.Enable) == 0 {
		analyzers := maps.Values(allAnalyzers)

		// To have a deterministic order.
		slices.SortFunc(analyzers, func(a *analysis.Analyzer, b *analysis.Analyzer) int {
			return strings.Compare(a.Name, b.Name)
		})

		return analyzers
	}

	var analyzers []*analysis.Analyzer
	for _, name := range uniqueNames(settings.Enable) {
		analyzers = append(analyzers, allAnalyzers[name])
	}

	return analyzers
}

func uniqueNames(names []string) []string {
	slices.Sort(names)
	return slices.Compact(names)
}
