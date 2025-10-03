package recovercheck

import (
	"github.com/cksidharthan/recovercheck"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.RecovercheckSettings) *goanalysis.Linter {
	var cfg recovercheck.RecovercheckSettings

	if settings != nil {
		cfg.SkipTestFiles = settings.SkipTestFiles
	}

	return goanalysis.
		NewLinterFromAnalyzer(recovercheck.New()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
