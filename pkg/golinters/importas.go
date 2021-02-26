package golinters

import (
	"fmt"

	"github.com/julz/importas" // nolint: misspell
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func NewImportAs(settings *config.ImportAsSettings) *goanalysis.Linter {
	analyzer := importas.Analyzer

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		if settings == nil {
			return
		}

		for alias, pkg := range *settings {
			err := analyzer.Flags.Set("alias", fmt.Sprintf("%s:%s", pkg, alias))
			if err != nil {
				lintCtx.Log.Errorf("failed to parse configuration: %v", err)
			}
		}
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
