package ineffassign

import (
	"strconv"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func New(settings *config.IneffAssignSettings) *goanalysis.Linter {
	analyzer := ineffassign.Analyzer

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithContextSetter(func(lintCtx *linter.Context) {
			if err := analyzer.Flags.Set("check-escaping-errors", strconv.FormatBool(settings.CheckEscapingErrors)); err != nil {
				lintCtx.Log.Errorf("failed to parse configuration: %v", err)
			}
		}).
		WithDesc("Detects when assignments to existing variables are not used").
		WithLoadMode(goanalysis.LoadModeSyntax)
}
