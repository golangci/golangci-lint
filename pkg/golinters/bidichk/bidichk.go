package bidichk

import (
	"strings"

	"github.com/breml/bidichk/pkg/bidichk"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.BiDiChkSettings) *goanalysis.Linter {
	a := bidichk.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		var opts []string

		if settings.LeftToRightEmbedding {
			opts = append(opts, "LEFT-TO-RIGHT-EMBEDDING")
		}
		if settings.RightToLeftEmbedding {
			opts = append(opts, "RIGHT-TO-LEFT-EMBEDDING")
		}
		if settings.PopDirectionalFormatting {
			opts = append(opts, "POP-DIRECTIONAL-FORMATTING")
		}
		if settings.LeftToRightOverride {
			opts = append(opts, "LEFT-TO-RIGHT-OVERRIDE")
		}
		if settings.RightToLeftOverride {
			opts = append(opts, "RIGHT-TO-LEFT-OVERRIDE")
		}
		if settings.LeftToRightIsolate {
			opts = append(opts, "LEFT-TO-RIGHT-ISOLATE")
		}
		if settings.RightToLeftIsolate {
			opts = append(opts, "RIGHT-TO-LEFT-ISOLATE")
		}
		if settings.FirstStrongIsolate {
			opts = append(opts, "FIRST-STRONG-ISOLATE")
		}
		if settings.PopDirectionalIsolate {
			opts = append(opts, "POP-DIRECTIONAL-ISOLATE")
		}

		cfg[a.Name] = map[string]any{
			"disallowed-runes": strings.Join(opts, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"Checks for dangerous unicode character sequences",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
