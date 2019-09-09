package lintersdb

import (
	"sort"
	"testing"

	"github.com/golangci/golangci-lint/pkg/golinters"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

//nolint:funlen
func TestGetEnabledLintersSet(t *testing.T) {
	type cs struct {
		cfg  config.Linters
		name string   // test case name
		def  []string // enabled by default linters
		exp  []string // alphabetically ordered enabled linter names
	}
	cases := []cs{
		{
			cfg: config.Linters{
				Disable: []string{golinters.MegacheckMetalinter{}.Name()},
			},
			name: "disable all linters from megacheck",
			def:  golinters.MegacheckMetalinter{}.DefaultChildLinterNames(),
			exp:  nil, // all disabled
		},
		{
			cfg: config.Linters{
				Disable: []string{golinters.MegacheckStaticcheckName},
			},
			name: "disable only staticcheck",
			def:  golinters.MegacheckMetalinter{}.DefaultChildLinterNames(),
			exp:  []string{golinters.MegacheckGosimpleName, golinters.MegacheckUnusedName},
		},
		{
			name: "don't merge into megacheck",
			def:  golinters.MegacheckMetalinter{}.DefaultChildLinterNames(),
			exp:  golinters.MegacheckMetalinter{}.DefaultChildLinterNames(),
		},
		{
			name: "expand megacheck",
			cfg: config.Linters{
				Enable: []string{golinters.MegacheckMetalinter{}.Name()},
			},
			def: nil,
			exp: golinters.MegacheckMetalinter{}.DefaultChildLinterNames(),
		},
		{
			name: "don't disable anything",
			def:  []string{"gofmt", "govet"},
			exp:  []string{"gofmt", "govet"},
		},
		{
			name: "enable gosec by gas alias",
			cfg: config.Linters{
				Enable: []string{"gas"},
			},
			exp: []string{"gosec"},
		},
		{
			name: "enable gosec by primary name",
			cfg: config.Linters{
				Enable: []string{"gosec"},
			},
			exp: []string{"gosec"},
		},
		{
			name: "enable gosec by both names",
			cfg: config.Linters{
				Enable: []string{"gosec", "gas"},
			},
			exp: []string{"gosec"},
		},
		{
			name: "disable gosec by gas alias",
			cfg: config.Linters{
				Disable: []string{"gas"},
			},
			def: []string{"gosec"},
		},
		{
			name: "disable gosec by primary name",
			cfg: config.Linters{
				Disable: []string{"gosec"},
			},
			def: []string{"gosec"},
		},
	}

	m := NewManager(nil)
	es := NewEnabledSet(m, NewValidator(m), nil, nil)
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var defaultLinters []*linter.Config
			for _, ln := range c.def {
				lc := m.GetLinterConfig(ln)
				assert.NotNil(t, lc, ln)
				defaultLinters = append(defaultLinters, lc)
			}

			els := es.build(&c.cfg, defaultLinters)
			var enabledLinters []string
			for ln, lc := range els {
				assert.Equal(t, ln, lc.Name())
				enabledLinters = append(enabledLinters, ln)
			}

			sort.Strings(enabledLinters)
			sort.Strings(c.exp)

			assert.Equal(t, c.exp, enabledLinters)
		})
	}
}
