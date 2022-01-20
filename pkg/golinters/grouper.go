package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	grouper "github.com/leonklingele/grouper/pkg/analyzer"
)

func NewGrouper(settings *config.GrouperSettings) *goanalysis.Linter {
	linterCfg := map[string]map[string]interface{}{}
	if settings != nil {
		linterCfg["grouper"] = map[string]interface{}{
			// const analyzer options
			"const-require-single-const": settings.ConstRequireSingleConst,
			"const-require-grouping":     settings.ConstRequireGrouping,
			// import analyzer options
			"import-require-single-import": settings.ImportRequireSingleImport,
			"import-require-grouping":      settings.ImportRequireGrouping,
			// type analyzer options
			"type-require-single-type": settings.TypeRequireSingleType,
			"type-require-grouping":    settings.TypeRequireGrouping,
			// var analyzer options
			"var-require-single-var": settings.VarRequireSingleVar,
			"var-require-grouping":   settings.VarRequireGrouping,
		}
	}

	return goanalysis.NewLinter(
		"grouper",
		"An analyzer to analyze expression groups.",
		[]*analysis.Analyzer{grouper.New()},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
