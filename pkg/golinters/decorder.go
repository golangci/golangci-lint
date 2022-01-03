package golinters

import (
	"strings"

	"gitlab.com/bosi/decorder"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewDecorder(settings *config.DecorderSettings) *goanalysis.Linter {
	a := decorder.Analyzer

	analyzers := []*analysis.Analyzer{a}

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			a.Name: {
				"dec-order":                     strings.Join(settings.DecOrder, ","),
				"disable-dec-num-check":         settings.DisableDecNumCheck,
				"disable-dec-order-check":       settings.DisableDecOrderCheck,
				"disable-init-func-first-check": settings.DisableInitFuncFirstCheck,
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		analyzers,
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
