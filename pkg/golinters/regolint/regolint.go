package regolint

import (
	"strings"

	"github.com/burdzwastaken/regolint/plugin"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.RegolintSettings, replacer *strings.Replacer) *goanalysis.Linter {
	cfg := make(map[string]any)

	if settings != nil {
		if settings.PolicyDir != "" {
			cfg["policy-dir"] = replacer.Replace(settings.PolicyDir)
		}
		if len(settings.PolicyFiles) > 0 {
			cfg["policy-files"] = settings.PolicyFiles
		}
		if len(settings.Disabled) > 0 {
			cfg["disabled"] = settings.Disabled
		}
		if len(settings.Exclude) > 0 {
			cfg["exclude"] = settings.Exclude
		}
	}

	p, err := plugin.New(cfg)
	if err != nil {
		return goanalysis.NewLinter(
			"regolint",
			"Define custom Go linting rules using Rego policies",
			nil,
			nil,
		)
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		return goanalysis.NewLinter(
			"regolint",
			"Define custom Go linting rules using Rego policies",
			nil,
			nil,
		)
	}

	return goanalysis.NewLinter(
		"regolint",
		"Define custom Go linting rules using Rego policies",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
