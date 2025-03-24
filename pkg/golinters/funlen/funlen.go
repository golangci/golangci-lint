package funlen

import (
	"github.com/ultraware/funlen"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

type Config struct {
	lineLimit      int
	stmtLimit      int
	ignoreComments bool
}

func New(settings *config.FunlenSettings) *goanalysis.Linter {
	cfg := Config{}
	if settings != nil {
		cfg.lineLimit = settings.Lines
		cfg.stmtLimit = settings.Statements
		cfg.ignoreComments = settings.IgnoreComments
	}

	a := funlen.NewAnalyzer(cfg.lineLimit, cfg.stmtLimit, cfg.ignoreComments)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
