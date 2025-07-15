package safebigint

import (
	"github.com/winder/safebigint"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.SafeBigIntSettings) *goanalysis.Linter {
	cfg := safebigint.LinterSettings{}
	if settings != nil {
		cfg.EnableTruncationCheck = !settings.DisableTruncationCheck
		cfg.EnableMutationCheck = !settings.DisableMutationCheck
	}

	analyzer, err := safebigint.NewAnalyzer(cfg)
	if err != nil {
		internal.LinterLogger.Fatalf("safebigint: new analyzer: %v", err)
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
