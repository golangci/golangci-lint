package golinters

import (
	"github.com/Ghvstcode/goBadWord/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

const gobadwordsname = "gobadwords"

func NewGoBadWords(settings *config.GoBadWordSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfg := make(map[string]map[string]interface{})
	if settings != nil {
		cfg[a.Name] = map[string]interface{}{
			"bad-words": settings.BadWords,
		}
	}

	return goanalysis.NewLinter(
		gobadwordsname,
		"Find occurrence of curse words or specified bad words",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}