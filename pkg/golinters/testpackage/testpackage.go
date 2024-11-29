package testpackage

import (
	"strings"

	"github.com/maratori/testpackage/pkg/testpackage"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TestpackageSettings) *goanalysis.Linter {
	a := testpackage.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				testpackage.SkipRegexpFlagName:    settings.SkipRegexp,
				testpackage.AllowPackagesFlagName: strings.Join(settings.AllowPackages, ","),
			},
		}
	}

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
