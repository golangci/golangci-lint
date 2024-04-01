package decorder

import (
	"strings"

	"gitlab.com/bosi/decorder"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.DecorderSettings) *goanalysis.Linter {
	a := decorder.Analyzer

	// disable all rules/checks by default
	cfg := map[string]any{
		"ignore-underscore-vars":        false,
		"disable-dec-num-check":         true,
		"disable-type-dec-num-check":    false,
		"disable-const-dec-num-check":   false,
		"disable-var-dec-num-check":     false,
		"disable-dec-order-check":       true,
		"disable-init-func-first-check": true,
	}

	if settings != nil {
		cfg["dec-order"] = strings.Join(settings.DecOrder, ",")
		cfg["ignore-underscore-vars"] = settings.IgnoreUnderscoreVars
		cfg["disable-dec-num-check"] = settings.DisableDecNumCheck
		cfg["disable-type-dec-num-check"] = settings.DisableTypeDecNumCheck
		cfg["disable-const-dec-num-check"] = settings.DisableConstDecNumCheck
		cfg["disable-var-dec-num-check"] = settings.DisableVarDecNumCheck
		cfg["disable-dec-order-check"] = settings.DisableDecOrderCheck
		cfg["disable-init-func-first-check"] = settings.DisableInitFuncFirstCheck
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		map[string]map[string]any{a.Name: cfg},
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
