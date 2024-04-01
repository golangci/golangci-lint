package grouper

import (
	grouper "github.com/leonklingele/grouper/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.GrouperSettings) *goanalysis.Linter {
	a := grouper.New()

	linterCfg := map[string]map[string]any{}
	if settings != nil {
		linterCfg[a.Name] = map[string]any{
			"const-require-single-const":   settings.ConstRequireSingleConst,
			"const-require-grouping":       settings.ConstRequireGrouping,
			"import-require-single-import": settings.ImportRequireSingleImport,
			"import-require-grouping":      settings.ImportRequireGrouping,
			"type-require-single-type":     settings.TypeRequireSingleType,
			"type-require-grouping":        settings.TypeRequireGrouping,
			"var-require-single-var":       settings.VarRequireSingleVar,
			"var-require-grouping":         settings.VarRequireGrouping,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"Analyze expression groups.",
		[]*analysis.Analyzer{a},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
