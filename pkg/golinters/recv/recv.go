package recv

import (
	"github.com/ldez/recv"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.RecvSettings) *goanalysis.Linter {
	cfg := recv.Config{}

	if settings != nil {
		cfg.MaxNameLength = settings.MaxNameLength
		cfg.TypeConsistency = settings.TypeConsistency
	}

	a := recv.New(cfg)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
