package golinters

import (
	"strings"

	"github.com/kulti/thelper/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewThelper(cfg *config.ThelperSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	opts := map[string]struct{}{
		"t_name":  {},
		"t_begin": {},
		"t_first": {},

		"f_name":  {},
		"f_begin": {},
		"f_first": {},

		"b_name":  {},
		"b_begin": {},
		"b_first": {},

		"tb_name":  {},
		"tb_begin": {},
		"tb_first": {},
	}

	if cfg != nil {
		applyTHelperOptions(cfg.Test, "t_", opts)
		applyTHelperOptions(cfg.Fuzz, "f_", opts)
		applyTHelperOptions(cfg.Benchmark, "b_", opts)
		applyTHelperOptions(cfg.TB, "tb_", opts)
	}

	if len(opts) == 0 {
		linterLogger.Fatalf("thelper: at least one option must be enabled")
	}

	var args []string
	for k := range opts {
		args = append(args, k)
	}

	cfgMap := map[string]map[string]interface{}{
		a.Name: {
			"checks": strings.Join(args, ","),
		},
	}

	return goanalysis.NewLinter(
		"thelper",
		"thelper detects Go test helpers without t.Helper() call and checks the consistency of test helpers",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func applyTHelperOptions(o config.ThelperOptions, prefix string, opts map[string]struct{}) {
	if o.Name != nil {
		if !*o.Name {
			delete(opts, prefix+"name")
		}
	}

	if o.Begin != nil {
		if !*o.Begin {
			delete(opts, prefix+"begin")
		}
	}

	if o.First != nil {
		if !*o.First {
			delete(opts, prefix+"first")
		}
	}
}
