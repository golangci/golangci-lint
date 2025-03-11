package dupword

import (
	"strings"

	"github.com/Abirdcfly/dupword"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"keyword": strings.Join(settings.Keywords, ","),
			"ignore":  strings.Join(settings.Ignore, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"checks for duplicate words in the source code",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
