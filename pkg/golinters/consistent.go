package golinters

import (
	"github.com/blizzy78/consistent"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewConsistent(settings *config.ConsistentSettings) *goanalysis.Linter {
	cfg := map[string]map[string]interface{}{}

	analyzer := consistent.NewAnalyzer()

	if settings != nil {
		valueMapping := map[string]string{
			"and-comp":     "andComp",
			"and-not":      "andNot",
			"compare-one":  "compareOne",
			"compare-zero": "compareZero",
			"equal-zero":   "equalZero",
		}

		analyzerCfg := map[string]interface{}{}

		set := func(key string, value string) {
			if v, ok := valueMapping[value]; ok {
				value = v
			}

			if value == "" {
				return
			}

			analyzerCfg[key] = value
		}

		set("params", settings.Params)
		set("returns", settings.Returns)
		set("typeParams", settings.TypeParams)
		set("singleImports", settings.SingleImports)
		set("newAllocs", settings.NewAllocs)
		set("makeAllocs", settings.MakeAllocs)
		set("hexLits", settings.HexLits)
		set("rangeChecks", settings.RangeChecks)
		set("andNOTs", settings.AndNOTs)
		set("floatLits", settings.FloatLits)
		set("lenChecks", settings.LenChecks)
		set("switchCases", settings.SwitchCases)
		set("switchDefaults", settings.SwitchDefaults)
		set("labelsRegexp", settings.LabelsRegexp)

		cfg[analyzer.Name] = analyzerCfg
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		"checks that common constructs are used consistently",
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
