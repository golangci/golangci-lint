package golinters

import (
	"github.com/jjti/go-spancheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSpancheck(settings *config.SpancheckSettings) *goanalysis.Linter {
	cfg := &spancheck.Config{}
	if settings != nil {
		cfg = &spancheck.Config{
			DisableEndCheck:                       settings.DisableEndCheck,
			EnableAll:                             settings.EnableAll,
			EnableRecordErrorCheck:                settings.EnableRecordErrorCheck,
			EnableSetStatusCheck:                  settings.EnableSetStatusCheck,
			IgnoreRecordErrorCheckSignaturesSlice: settings.IgnoreRecordErrorCheckSignatures,
			IgnoreSetStatusCheckSignaturesSlice:   settings.IgnoreSetStatusCheckSignatures,
		}
	}

	a := spancheck.NewAnalyzerWithConfig(cfg)

	return goanalysis.
		NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
