package golinters

import (
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/butuzov/ireturn/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewIreturn(settings *config.IreturnSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		ireturnSettings(a.Name, settings),
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func ireturnSettings(name string, s *config.IreturnSettings) map[string]map[string]interface{} {
	if s == nil {
		return nil
	}

	return map[string]map[string]interface{}{
		name: {
			"allow":  strings.Join(s.Allow, ","),
			"reject": strings.Join(s.Reject, ","),
		},
	}
}
