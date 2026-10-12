package predeclared

import (
	"strings"

	"github.com/nishanths/predeclared/passes/predeclared"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.PredeclaredSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"ignore":   strings.Join(settings.Ignore, ","),
			"q":        settings.Qualified,
			"pkglevel": false,
			"mode":     modeToString(settings.Mode),
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(predeclared.Analyzer).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func modeToString(cfg []string) string {
	var mode strings.Builder

	for _, m := range cfg {
		switch m {
		case "declare":
			mode.WriteString("d")
		case "shadow":
			mode.WriteString("s")
		default:
			internal.LinterLogger.Warnf("predeclared: unsupported mode %q", m)
		}
	}

	if mode.Len() == 0 {
		return "s"
	}

	return mode.String()
}
