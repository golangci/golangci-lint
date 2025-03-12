package reassign

import (
	"fmt"
	"strings"

	"github.com/curioswitch/go-reassign"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.ReassignSettings) *goanalysis.Linter {
	a := reassign.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil && len(settings.Patterns) > 0 {
		cfg = map[string]map[string]any{
			a.Name: {
				reassign.FlagPattern: fmt.Sprintf("^(%s)$", strings.Join(settings.Patterns, "|")),
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
