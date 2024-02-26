package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type dummyLogger struct{}

func (d dummyLogger) Fatalf(_ string, _ ...any) {}

func (d dummyLogger) Panicf(_ string, _ ...any) {}

func (d dummyLogger) Errorf(_ string, _ ...any) {}

func (d dummyLogger) Warnf(_ string, _ ...any) {}

func (d dummyLogger) Infof(_ string, _ ...any) {}

func (d dummyLogger) Child(_ string) logutils.Log {
	return nil
}

func (d dummyLogger) SetLevel(_ logutils.LogLevel) {}

func TestEnabledSet_GetEnabledLintersMap(t *testing.T) {
	m := NewManager(nil, nil)

	cfg := config.NewDefault()

	cfg.Linters.DisableAll = true
	cfg.Linters.Enable = []string{"gofmt"}

	es := NewEnabledSet(m, NewValidator(m), dummyLogger{}, cfg)

	lintersMap, err := es.GetEnabledLintersMap()
	require.NoError(t, err)

	gofmtConfigs := m.GetLinterConfigs("gofmt")
	typecheckConfigs := m.GetLinterConfigs("typecheck")

	expected := map[string]*linter.Config{
		"gofmt":     gofmtConfigs[0],
		"typecheck": typecheckConfigs[0],
	}

	assert.Equal(t, expected, lintersMap)
}

func TestEnabledSet_GetOptimizedLinters(t *testing.T) {
	m := NewManager(nil, nil)

	cfg := config.NewDefault()

	cfg.Linters.DisableAll = true
	cfg.Linters.Enable = []string{"gofmt"}

	es := NewEnabledSet(m, NewValidator(m), dummyLogger{}, cfg)

	optimizedLinters, err := es.GetOptimizedLinters()
	require.NoError(t, err)

	gofmtConfigs := m.GetLinterConfigs("gofmt")
	typecheckConfigs := m.GetLinterConfigs("typecheck")

	var gaLinters []*goanalysis.Linter
	for _, l := range gofmtConfigs {
		gaLinters = append(gaLinters, l.Linter.(*goanalysis.Linter))
	}
	for _, l := range typecheckConfigs {
		gaLinters = append(gaLinters, l.Linter.(*goanalysis.Linter))
	}

	mlConfig := &linter.Config{
		Linter:    goanalysis.NewMetaLinter(gaLinters),
		InPresets: []string{"bugs", "format"},
	}

	expected := []*linter.Config{mlConfig.WithLoadForGoAnalysis()}

	assert.Equal(t, expected, optimizedLinters)
}

func TestEnabledSet_build(t *testing.T) {
	type cs struct {
		cfg  config.Linters
		name string   // test case name
		def  []string // enabled by default linters
		exp  []string // alphabetically ordered enabled linter names
	}

	allMegacheckLinterNames := []string{"gosimple", "staticcheck", "unused"}

	cases := []cs{
		{
			cfg: config.Linters{
				Disable: []string{"megacheck"},
			},
			name: "disable all linters from megacheck",
			def:  allMegacheckLinterNames,
			exp:  []string{"typecheck"}, // all disabled
		},
		{
			cfg: config.Linters{
				Disable: []string{"staticcheck"},
			},
			name: "disable only staticcheck",
			def:  allMegacheckLinterNames,
			exp:  []string{"gosimple", "typecheck", "unused"},
		},
		{
			name: "don't merge into megacheck",
			def:  allMegacheckLinterNames,
			exp:  []string{"gosimple", "staticcheck", "typecheck", "unused"},
		},
		{
			name: "expand megacheck",
			cfg: config.Linters{
				Enable: []string{"megacheck"},
			},
			def: nil,
			exp: []string{"gosimple", "staticcheck", "typecheck", "unused"},
		},
		{
			name: "don't disable anything",
			def:  []string{"gofmt", "govet", "typecheck"},
			exp:  []string{"gofmt", "govet", "typecheck"},
		},
		{
			name: "enable gosec by gas alias",
			cfg: config.Linters{
				Enable: []string{"gas"},
			},
			exp: []string{"gosec", "typecheck"},
		},
		{
			name: "enable gosec by primary name",
			cfg: config.Linters{
				Enable: []string{"gosec"},
			},
			exp: []string{"gosec", "typecheck"},
		},
		{
			name: "enable gosec by both names",
			cfg: config.Linters{
				Enable: []string{"gosec", "gas"},
			},
			exp: []string{"gosec", "typecheck"},
		},
		{
			name: "disable gosec by gas alias",
			cfg: config.Linters{
				Disable: []string{"gas"},
			},
			def: []string{"gosec"},
			exp: []string{"typecheck"},
		},
		{
			name: "disable gosec by primary name",
			cfg: config.Linters{
				Disable: []string{"gosec"},
			},
			def: []string{"gosec"},
			exp: []string{"typecheck"},
		},
	}

	m := NewManager(nil, nil)
	es := NewEnabledSet(m, NewValidator(m), dummyLogger{}, nil)

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var defaultLinters []*linter.Config
			for _, ln := range c.def {
				lcs := m.GetLinterConfigs(ln)
				assert.NotNil(t, lcs, ln)
				defaultLinters = append(defaultLinters, lcs...)
			}

			els := es.build(&c.cfg, defaultLinters)
			var enabledLinters []string
			for ln, lc := range els {
				assert.Equal(t, ln, lc.Name())
				enabledLinters = append(enabledLinters, ln)
			}

			assert.ElementsMatch(t, c.exp, enabledLinters)
		})
	}
}

func TestEnabledSet_combineGoAnalysisLinters(t *testing.T) {
	m := NewManager(nil, nil)

	es := NewEnabledSet(m, NewValidator(m), dummyLogger{}, config.NewDefault())

	foo := goanalysis.NewLinter("foo", "example foo", nil, nil).WithLoadMode(goanalysis.LoadModeTypesInfo)
	bar := goanalysis.NewLinter("bar", "example bar", nil, nil).WithLoadMode(goanalysis.LoadModeTypesInfo)

	testCases := []struct {
		desc     string
		linters  map[string]*linter.Config
		expected map[string]*linter.Config
	}{
		{
			desc: "no combined, one linter",
			linters: map[string]*linter.Config{
				"foo": {
					Linter:    foo,
					InPresets: []string{"A"},
				},
			},
			expected: map[string]*linter.Config{
				"foo": {
					Linter:    foo,
					InPresets: []string{"A"},
				},
			},
		},
		{
			desc: "combined, several linters",
			linters: map[string]*linter.Config{
				"foo": {
					Linter:    foo,
					InPresets: []string{"A"},
				},
				"bar": {
					Linter:    bar,
					InPresets: []string{"B"},
				},
			},
			expected: func() map[string]*linter.Config {
				mlConfig := &linter.Config{
					Linter:    goanalysis.NewMetaLinter([]*goanalysis.Linter{bar, foo}),
					InPresets: []string{"A", "B"},
				}

				return map[string]*linter.Config{
					"goanalysis_metalinter": mlConfig.WithLoadForGoAnalysis(),
				}
			}(),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			es.combineGoAnalysisLinters(test.linters)

			assert.Equal(t, test.expected, test.linters)
		})
	}
}
