package golinters

import (
	"github.com/drichelson/gokeyword"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

const (
	goKeywordName        = "gokeyword"
	goKeywordDescription = "detects presence of the go keyword"
)

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {
	a := gokeyword.New()

	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil {
		cfgMap[a.Name] = map[string]interface{}{
			"details": cfg.Details,
		}
	}

	return goanalysis.NewLinter(
		goKeywordName,
		goKeywordDescription,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
