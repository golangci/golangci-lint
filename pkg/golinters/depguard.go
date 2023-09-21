package golinters

import (
	"os"

	"github.com/OpenPeeDeeP/depguard/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const modFileName string = "go.mod"

func NewDepguard(settings *config.DepGuardSettings) *goanalysis.Linter {
	conf := depguard.LinterSettings{}

	if currentModule := readCurrentModule(); currentModule != "" {
		conf[currentModule] = &depguard.List{Allow: []string{currentModule}}
	}

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

func readCurrentModule() string {
	data, err := os.ReadFile(modFileName)
	if err != nil {
		return ""
	}

	modFile, err := modfile.Parse(modFileName, data, nil)
	if err != nil {
		return ""
	}
	return modFile.Module.Mod.String()
}
