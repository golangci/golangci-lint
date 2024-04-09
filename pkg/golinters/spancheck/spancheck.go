package spancheck

import (
	"github.com/jjti/go-spancheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.SpancheckSettings) *goanalysis.Linter {
	cfg := spancheck.NewDefaultConfig()

	if settings != nil {
		if settings.Checks != nil {
			cfg.EnabledChecks = settings.Checks
		}

		if settings.IgnoreCheckSignatures != nil {
			cfg.IgnoreChecksSignaturesSlice = settings.IgnoreCheckSignatures
		}
	}

	a := spancheck.NewAnalyzerWithConfig(cfg)

	return goanalysis.
		NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
