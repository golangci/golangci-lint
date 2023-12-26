package golinters

import (
	"honnef.co/go/tools/staticcheck"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewStaticcheck(settings *config.CommonStaticCheckSettings) *goanalysis.Linter {
	cfg := commonStaticCheckConfig(settings)
	analyzers := setupStaticCheckAnalyzers(staticcheck.Analyzers, settings.GetGoVersion(), cfg.Checks)

	return goanalysis.NewLinter(
		"staticcheck",
		"It's a set of rules from staticcheck. It's not the same thing as the staticcheck binary."+
			" The author of staticcheck doesn't support or approve the use of staticcheck as a library inside golangci-lint.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
