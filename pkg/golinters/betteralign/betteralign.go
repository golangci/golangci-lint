package betteralign

import (
	"strconv"
	"strings"

	"github.com/dkorunic/betteralign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func New(settings *config.BetteralignSettings) *goanalysis.Linter {
	a := betteralign.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithContextSetter(func(ctx *linter.Context) {
		if settings == nil {
			return
		}

		if err := a.Flags.Set("test_files", strconv.FormatBool(settings.TestFiles)); err != nil {
			ctx.Log.Infof("failed to parse configuration: %v", err)
		}

		if err := a.Flags.Set("generated_files", strconv.FormatBool(settings.GeneratedFiles)); err != nil {
			ctx.Log.Infof("failed to parse configuration: %v", err)
		}

		if len(settings.ExcludeFiles) > 0 {
			if err := a.Flags.Set("exclude_files", strings.Join(settings.ExcludeFiles, ",")); err != nil {
				ctx.Log.Infof("failed to parse configuration: %v", err)
			}
		}

		if len(settings.ExcludeDirs) > 0 {
			if err := a.Flags.Set("exclude_dirs", strings.Join(settings.ExcludeDirs, ",")); err != nil {
				ctx.Log.Infof("failed to parse configuration: %v", err)
			}
		}
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
