package thelper

import (
	"maps"
	"slices"
	"strings"

	"github.com/kulti/thelper/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

func New(settings *config.ThelperSettings) *goanalysis.Linter {
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

	if settings != nil {
		applyTHelperOptions(settings.Test, "t_", opts)
		applyTHelperOptions(settings.Fuzz, "f_", opts)
		applyTHelperOptions(settings.Benchmark, "b_", opts)
		applyTHelperOptions(settings.TB, "tb_", opts)
	}

	if len(opts) == 0 {
		internal.LinterLogger.Fatalf("thelper: at least one option must be enabled")
	}

	args := slices.Collect(maps.Keys(opts))

	cfg := map[string]map[string]any{
		a.Name: {
			"checks": strings.Join(args, ","),
		},
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
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
