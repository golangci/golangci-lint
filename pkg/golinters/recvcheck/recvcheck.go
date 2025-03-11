package recvcheck

import (
	"github.com/raeperd/recvcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.RecvcheckSettings) *goanalysis.Linter {
	var cfg recvcheck.Settings

	if settings != nil {
		cfg.DisableBuiltin = settings.DisableBuiltin
		cfg.Exclusions = settings.Exclusions
	}

	a := recvcheck.NewAnalyzer(cfg)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
