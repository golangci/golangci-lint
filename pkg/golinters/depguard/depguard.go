package depguard

import (
	"strings"

	"github.com/OpenPeeDeeP/depguard/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func New(settings *config.DepGuardSettings, replacer *strings.Replacer) *goanalysis.Linter {
	conf := depguard.LinterSettings{}

	if settings != nil {
		for s, rule := range settings.Rules {
			var extendedPatterns []string
			for _, file := range rule.Files {
				extendedPatterns = append(extendedPatterns, replacer.Replace(file))
			}

			list := &depguard.List{
				ListMode: rule.ListMode,
				Files:    extendedPatterns,
				Allow:    rule.Allow,
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

	a := depguard.NewUncompiledAnalyzer(&conf)

	return goanalysis.NewLinter(
		a.Analyzer.Name,
		a.Analyzer.Doc,
		[]*analysis.Analyzer{a.Analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		err := a.Compile()
		if err != nil {
			lintCtx.Log.Errorf("create analyzer: %v", err)
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
