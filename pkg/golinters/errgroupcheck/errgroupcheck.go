package errgroupcheck

import (
	"github.com/alexbagnolini/errgroupcheck"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func New(cfg *config.ErrGroupCheckSettings) *goanalysis.Linter {
	var setts = errgroupcheck.DefaultSettings()

	if cfg != nil {
		setts.RequireWait = cfg.RequireWait
	}

	cfgMap := map[string]map[string]any{}

	analyzer := errgroupcheck.NewAnalyzer(setts)

	if cfg != nil {
		cfgMap[analyzer.Name] = map[string]any{
			"require-wait": cfg.RequireWait,
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
