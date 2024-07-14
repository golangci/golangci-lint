package iface

import (
	"slices"

	"github.com/uudashr/iface/duplicate"
	"github.com/uudashr/iface/empty"
	"github.com/uudashr/iface/opaque"
	"github.com/uudashr/iface/unused"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

var allAnalyzers = []*analysis.Analyzer{
	unused.Analyzer,
	empty.Analyzer,
	duplicate.Analyzer,
	opaque.Analyzer,
}

func New(settings *config.IfaceSettings) *goanalysis.Linter {
	var conf map[string]map[string]any

	analyzers := analyzersFromSettings(settings)

	return goanalysis.NewLinter(
		"iface",
		"Detect the incorrect use of interfaces, helping developers avoid interface pollution.",
		analyzers,
		conf,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func analyzersFromSettings(settings *config.IfaceSettings) []*analysis.Analyzer {
	if settings == nil || len(settings.Enable) == 0 {
		return allAnalyzers
	}

	enabledNames := uniqueNames(settings.Enable)

	var analyzers []*analysis.Analyzer

	for _, a := range allAnalyzers {
		found := slices.ContainsFunc(enabledNames, func(name string) bool {
			return name == a.Name
		})

		if !found {
			continue
		}

		analyzers = append(analyzers, a)
	}

	return analyzers
}

func uniqueNames(names []string) []string {
	if len(names) == 0 {
		return nil
	}

	namesMap := map[string]struct{}{}
	for _, name := range names {
		namesMap[name] = struct{}{}
	}

	uniqueNames := make([]string, 0, len(namesMap))

	for name := range namesMap {
		uniqueNames = append(uniqueNames, name)
	}
	return uniqueNames
}
