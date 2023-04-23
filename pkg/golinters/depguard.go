package golinters

import (
	"github.com/OpenPeeDeeP/depguard/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewDepguard(settings *config.DepGuardSettings) *goanalysis.Linter {
	conf := depguard.LinterSettings{}

	if settings != nil {
		for s, rule := range settings.Rules {
			list := &depguard.List{
				Files: rule.Files,
				Allow: rule.Allow,
			}

			// because of bug with Viper parsing (split on dot) we use a list of struct instead of a map.
			// https://github.com/spf13/viper/issues/324
			// https://github.com/golangci/golangci-lint/issues/3749#issuecomment-1492536630

			deny := map[string]string{}
			for _, r := range rule.Deny {
				deny[r.Pkg] = r.Desc
			}
			list.Deny = deny

			conf[s] = list
		}
	}

	a, err := depguard.NewAnalyzer(&conf)
	if err != nil {
		linterLogger.Fatalf("depguard: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
