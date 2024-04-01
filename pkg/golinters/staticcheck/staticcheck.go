package staticcheck

import (
	"honnef.co/go/tools/staticcheck"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

func New(settings *config.StaticCheckSettings) *goanalysis.Linter {
	cfg := internal.StaticCheckConfig(settings)
	analyzers := internal.SetupStaticCheckAnalyzers(staticcheck.Analyzers, internal.GetGoVersion(settings), cfg.Checks)

	return goanalysis.NewLinter(
		"staticcheck",
		"It's a set of rules from staticcheck. It's not the same thing as the staticcheck binary."+
			" The author of staticcheck doesn't support or approve the use of staticcheck as a library inside golangci-lint.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
