package lintersdb

import (
	"sort"
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/stretchr/testify/assert"
)

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
				Disable: []string{"megacheck"},
			},
			name: "disable all linters from megacheck",
			def:  getAllMegacheckSubLinterNames(),
		},
		{
			cfg: config.Linters{
				Disable: []string{"staticcheck"},
			},
			name: "disable only staticcheck",
			def:  getAllMegacheckSubLinterNames(),
			exp:  []string{"megacheck.{unused,gosimple}"},
		},
		{
			name: "merge into megacheck",
			def:  getAllMegacheckSubLinterNames(),
			exp:  []string{"megacheck"},
		},
		{
			name: "don't disable anything",
			def:  []string{"gofmt", "govet"},
			exp:  []string{"gofmt", "govet"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defaultLinters := []linter.Config{}
			for _, ln := range c.def {
				defaultLinters = append(defaultLinters, *getLinterConfig(ln))
			}
			els := getEnabledLintersSet(&c.cfg, defaultLinters)
			var enabledLinters []string
			for ln, lc := range els {
				assert.Equal(t, ln, lc.Linter.Name())
				enabledLinters = append(enabledLinters, ln)
			}

			sort.Strings(enabledLinters)
			sort.Strings(c.exp)

			assert.Equal(t, c.exp, enabledLinters)
		})
	}
}
