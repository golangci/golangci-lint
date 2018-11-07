package test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/test/testshared"
)

func inSlice(s []string, v string) bool {
	for _, sv := range s {
		if sv == v {
			return true
		}
	}

	return false
}

func getEnabledByDefaultFastLintersExcept(except ...string) []string {
	m := lintersdb.NewManager()
	ebdl := m.GetAllEnabledByDefaultLinters()
	ret := []string{}
	for _, lc := range ebdl {
		if lc.NeedsSSARepr {
			continue
		}

		if !inSlice(except, lc.Name()) {
			ret = append(ret, lc.Name())
		}
	}

	return ret
}

func getAllFastLintersWith(with ...string) []string {
	linters := lintersdb.NewManager().GetAllSupportedLinterConfigs()
	ret := append([]string{}, with...)
	for _, lc := range linters {
		if lc.NeedsSSARepr {
			continue
		}
		ret = append(ret, lc.Name())
	}

	return ret
}

func getEnabledByDefaultLinters() []string {
	ebdl := lintersdb.NewManager().GetAllEnabledByDefaultLinters()
	ret := []string{}
	for _, lc := range ebdl {
		ret = append(ret, lc.Name())
	}

	return ret
}

func getEnabledByDefaultFastLintersWith(with ...string) []string {
	ebdl := lintersdb.NewManager().GetAllEnabledByDefaultLinters()
	ret := append([]string{}, with...)
	for _, lc := range ebdl {
		if lc.NeedsSSARepr {
			continue
		}

		ret = append(ret, lc.Name())
	}

	return ret
}

func mergeMegacheck(linters []string) []string {
	if inSlice(linters, "staticcheck") &&
		inSlice(linters, "gosimple") &&
		inSlice(linters, "unused") {
		ret := []string{"megacheck"}
		for _, linter := range linters {
			if !inSlice([]string{"staticcheck", "gosimple", "unused"}, linter) {
				ret = append(ret, linter)
			}
		}

		return ret
	}

	return linters
}

func TestEnabledLinters(t *testing.T) {
	type tc struct {
		name           string
		cfg            string
		el             []string
		args           string
		noImplicitFast bool
	}

	cases := []tc{
		{
			name: "disable govet in config",
			cfg: `
			linters:
				disable:
					- govet
			`,
			el: getEnabledByDefaultFastLintersExcept("govet"),
		},
		{
			name: "enable golint in config",
			cfg: `
			linters:
				enable:
					- golint
			`,
			el: getEnabledByDefaultFastLintersWith("golint"),
		},
		{
			name: "disable govet in cmd",
			args: "-Dgovet",
			el:   getEnabledByDefaultFastLintersExcept("govet"),
		},
		{
			name: "enable gofmt in cmd and enable golint in config",
			args: "-Egofmt",
			cfg: `
			linters:
				enable:
					- golint
			`,
			el: getEnabledByDefaultFastLintersWith("golint", "gofmt"),
		},
		{
			name: "fast option in config",
			cfg: `
			linters:
				fast: true
			`,
			el:             getEnabledByDefaultFastLintersWith(),
			noImplicitFast: true,
		},
		{
			name: "explicitly unset fast option in config",
			cfg: `
			linters:
				fast: false
			`,
			el:             getEnabledByDefaultLinters(),
			noImplicitFast: true,
		},
		{
			name:           "set fast option in command-line",
			args:           "--fast",
			el:             getEnabledByDefaultFastLintersWith(),
			noImplicitFast: true,
		},
		{
			name: "fast option in command-line has higher priority to enable",
			cfg: `
			linters:
				fast: false
			`,
			args:           "--fast",
			el:             getEnabledByDefaultFastLintersWith(),
			noImplicitFast: true,
		},
		{
			name: "fast option in command-line has higher priority to disable",
			cfg: `
			linters:
				fast: true
			`,
			args:           "--fast=false",
			el:             getEnabledByDefaultLinters(),
			noImplicitFast: true,
		},
		{
			name:           "fast option combined with enable and enable-all",
			args:           "--enable-all --fast --enable=megacheck",
			el:             getAllFastLintersWith("megacheck"),
			noImplicitFast: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			runArgs := []string{"-v"}
			if !c.noImplicitFast {
				runArgs = append(runArgs, "--fast")
			}
			if c.args != "" {
				runArgs = append(runArgs, strings.Split(c.args, " ")...)
			}
			r := testshared.NewLintRunner(t).RunWithYamlConfig(c.cfg, runArgs...)
			el := mergeMegacheck(c.el)
			sort.StringSlice(el).Sort()

			expectedLine := fmt.Sprintf("Active %d linters: [%s]", len(el), strings.Join(el, " "))
			r.ExpectOutputContains(expectedLine)
		})
	}
}
