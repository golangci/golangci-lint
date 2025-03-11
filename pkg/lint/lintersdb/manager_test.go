package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func TestManager_GetEnabledLintersMap(t *testing.T) {
	cfg := config.NewDefault()
	cfg.Linters.Default = config.GroupNone
	cfg.Linters.Enable = []string{"gofmt"}

	m, err := NewManager(logutils.NewStderrLog("skip"), cfg, NewLinterBuilder())
	require.NoError(t, err)

	lintersMap, err := m.GetEnabledLintersMap()
	require.NoError(t, err)

	gofmtConfigs := m.GetLinterConfigs("gofmt")
	typecheckConfigs := m.GetLinterConfigs("typecheck")

	expected := map[string]*linter.Config{
		"gofmt":     gofmtConfigs[0],
		"typecheck": typecheckConfigs[0],
	}

	assert.Equal(t, expected, lintersMap)
}

func TestManager_GetOptimizedLinters(t *testing.T) {
	cfg := config.NewDefault()
	cfg.Linters.Default = config.GroupNone
	cfg.Linters.Enable = []string{"gofmt"}

	m, err := NewManager(logutils.NewStderrLog("skip"), cfg, NewLinterBuilder())
	require.NoError(t, err)

	optimizedLinters, err := m.GetOptimizedLinters()
	require.NoError(t, err)

	var gaLinters []*goanalysis.Linter
	for _, l := range m.GetLinterConfigs("gofmt") {
		gaLinters = append(gaLinters, l.Linter.(*goanalysis.Linter))
	}
	for _, l := range m.GetLinterConfigs("typecheck") {
		gaLinters = append(gaLinters, l.Linter.(*goanalysis.Linter))
	}

	mlConfig := &linter.Config{
		Linter: goanalysis.NewMetaLinter(gaLinters),
	}

	expected := []*linter.Config{mlConfig.WithLoadFiles()}

	assert.Equal(t, expected, optimizedLinters)
}

func TestManager_build(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      *config.Config
		expected []string // alphabetically ordered enabled linter names
	}{
		{
			desc: "enable gosec by primary name",
			cfg: &config.Config{
				Linters: config.Linters{
					Default: config.GroupNone,
					Enable:  []string{"gosec"},
				},
			},
			expected: []string{"gosec", "typecheck"},
		},
		{
			desc: "disable gosec by primary name",
			cfg: &config.Config{
				Linters: config.Linters{
					Default: config.GroupNone,
					Disable: []string{"gosec"},
				},
			},
			expected: []string{"typecheck"},
		},
		{
			desc: "linters and formatters",
			cfg: &config.Config{
				Linters: config.Linters{
					Default: config.GroupNone,
					Enable:  []string{"gosec"},
				},
				Formatters: config.Formatters{
					Enable: []string{"gofmt"},
				},
			},
			expected: []string{"gosec", "gofmt", "typecheck"},
		},
		{
			desc: "linters and formatters but linters configuration disables the formatter",
			cfg: &config.Config{
				Linters: config.Linters{
					Default: config.GroupNone,
					Enable:  []string{"gosec"},
					Disable: []string{"gofmt"},
				},
				Formatters: config.Formatters{
					Enable: []string{"gofmt"},
				},
			},
			expected: []string{"gosec", "typecheck"},
		},
		{
			desc: "only formatters",
			cfg: &config.Config{
				Linters: config.Linters{
					Default: config.GroupNone,
				},
				Formatters: config.Formatters{
					Enable: []string{"gofmt"},
				},
			},
			expected: []string{"gofmt", "typecheck"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			m, err := NewManager(logutils.NewStderrLog("skip"), test.cfg, NewLinterBuilder())
			require.NoError(t, err)

			els, err := m.build()
			require.NoError(t, err)

			var enabledLinters []string
			for ln, lc := range els {
				assert.Equal(t, ln, lc.Name())
				enabledLinters = append(enabledLinters, ln)
			}

			assert.ElementsMatch(t, test.expected, enabledLinters)
		})
	}
}

func TestManager_combineGoAnalysisLinters(t *testing.T) {
	m, err := NewManager(nil, nil)
	require.NoError(t, err)

	fooTyped := goanalysis.NewLinter("foo", "example foo", nil, nil).WithLoadMode(goanalysis.LoadModeTypesInfo)
	barTyped := goanalysis.NewLinter("bar", "example bar", nil, nil).WithLoadMode(goanalysis.LoadModeTypesInfo)

	fooSyntax := goanalysis.NewLinter("foo", "example foo", nil, nil).WithLoadMode(goanalysis.LoadModeSyntax)
	barSyntax := goanalysis.NewLinter("bar", "example bar", nil, nil).WithLoadMode(goanalysis.LoadModeSyntax)

	testCases := []struct {
		desc     string
		linters  map[string]*linter.Config
		expected map[string]*linter.Config
	}{
		{
			desc: "no combined, one linter",
			linters: map[string]*linter.Config{
				"foo": {
					Linter: fooTyped,
				},
			},
			expected: map[string]*linter.Config{
				"foo": {
					Linter: fooTyped,
				},
			},
		},
		{
			desc: "combined, several linters (typed)",
			linters: map[string]*linter.Config{
				"foo": {
					Linter: fooTyped,
				},
				"bar": {
					Linter: barTyped,
				},
			},
			expected: func() map[string]*linter.Config {
				mlConfig := &linter.Config{
					Linter: goanalysis.NewMetaLinter([]*goanalysis.Linter{barTyped, fooTyped}),
				}

				return map[string]*linter.Config{
					"goanalysis_metalinter": mlConfig,
				}
			}(),
		},
		{
			desc: "combined, several linters (different LoadMode)",
			linters: map[string]*linter.Config{
				"foo": {
					Linter:   fooTyped,
					LoadMode: packages.NeedName,
				},
				"bar": {
					Linter:   barTyped,
					LoadMode: packages.NeedTypesSizes,
				},
			},
			expected: func() map[string]*linter.Config {
				mlConfig := &linter.Config{
					Linter:   goanalysis.NewMetaLinter([]*goanalysis.Linter{barTyped, fooTyped}),
					LoadMode: packages.NeedName | packages.NeedTypesSizes,
				}

				return map[string]*linter.Config{
					"goanalysis_metalinter": mlConfig,
				}
			}(),
		},
		{
			desc: "combined, several linters (same LoadMode)",
			linters: map[string]*linter.Config{
				"foo": {
					Linter:   fooTyped,
					LoadMode: packages.NeedName,
				},
				"bar": {
					Linter:   barTyped,
					LoadMode: packages.NeedName,
				},
			},
			expected: func() map[string]*linter.Config {
				mlConfig := &linter.Config{
					Linter:   goanalysis.NewMetaLinter([]*goanalysis.Linter{barTyped, fooTyped}),
					LoadMode: packages.NeedName,
				}

				return map[string]*linter.Config{
					"goanalysis_metalinter": mlConfig,
				}
			}(),
		},
		{
			desc: "combined, several linters (syntax)",
			linters: map[string]*linter.Config{
				"foo": {
					Linter: fooSyntax,
				},
				"bar": {
					Linter: barSyntax,
				},
			},
			expected: func() map[string]*linter.Config {
				mlConfig := &linter.Config{
					Linter: goanalysis.NewMetaLinter([]*goanalysis.Linter{barSyntax, fooSyntax}),
				}

				return map[string]*linter.Config{
					"goanalysis_metalinter": mlConfig,
				}
			}(),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			m.combineGoAnalysisLinters(test.linters)

			assert.Equal(t, test.expected, test.linters)
		})
	}
}
