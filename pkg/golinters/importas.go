package golinters

import (
	"fmt"
	"strconv"

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

		err := analyzer.Flags.Set("no-unaliased", strconv.FormatBool(settings.NoUnaliased))
		if err != nil {
			lintCtx.Log.Errorf("failed to parse configuration: %v", err)
		}

		for alias, pkg := range settings.Alias {
			err := analyzer.Flags.Set("alias", fmt.Sprintf("%s:%s", pkg, alias))
			if err != nil {
				lintCtx.Log.Errorf("failed to parse configuration: %v", err)
			}
		}
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
