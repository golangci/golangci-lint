package ireturn

import (
	"strings"

	"github.com/butuzov/ireturn/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.IreturnSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"allow":    strings.Join(settings.Allow, ","),
			"reject":   strings.Join(settings.Reject, ","),
			"nonolint": true,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
