package recovercheck

import (
	"github.com/cksidharthan/recovercheck"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(recovercheck.New()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
