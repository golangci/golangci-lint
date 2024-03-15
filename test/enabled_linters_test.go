package test

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/test/testshared"
)

func TestEnabledLinters(t *testing.T) {
	// require to display the message "Active x linters: [x,y]"
	t.Setenv(logutils.EnvTestRun, "1")

	cases := []struct {
		name           string
		cfg            string
		enabledLinters []string
		args           []string
		noImplicitFast bool
	}{
		{
			name: "disable govet in config",
			cfg: `
			linters:
				disable:
					- govet
			`,
			enabledLinters: getEnabledByDefaultFastLintersExcept(t, "govet"),
		},
		{
			name: "enable revive in config",
			cfg: `
			linters:
				enable:
					- revive
			`,
			enabledLinters: getEnabledByDefaultFastLintersWith(t, "revive"),
		},
		{
			name:           "disable govet in cmd",
			args:           []string{"-Dgovet"},
			enabledLinters: getEnabledByDefaultFastLintersExcept(t, "govet"),
		},
		{
			name: "enable gofmt in cmd and enable revive in config",
			args: []string{"-Egofmt"},
			cfg: `
			linters:
				enable:
					- revive
			`,
			enabledLinters: getEnabledByDefaultFastLintersWith(t, "revive", "gofmt"),
		},
		{
			name: "fast option in config",
			cfg: `
			linters:
				fast: true
			`,
			enabledLinters: getEnabledByDefaultFastLintersWith(t),
			noImplicitFast: true,
		},
		{
			name: "explicitly unset fast option in config",
			cfg: `
			linters:
				fast: false
			`,
			enabledLinters: getEnabledByDefaultLinters(t),
			noImplicitFast: true,
		},
		{
			name:           "set fast option in command-line",
			args:           []string{"--fast"},
			enabledLinters: getEnabledByDefaultFastLintersWith(t),
			noImplicitFast: true,
		},
		{
			name: "fast option in command-line has higher priority to enable",
			cfg: `
			linters:
				fast: false
			`,
			args:           []string{"--fast"},
			enabledLinters: getEnabledByDefaultFastLintersWith(t),
			noImplicitFast: true,
		},
		{
			name: "fast option in command-line has higher priority to disable",
			cfg: `
			linters:
				fast: true
			`,
			args:           []string{"--fast=false"},
			enabledLinters: getEnabledByDefaultLinters(t),
			noImplicitFast: true,
		},
		{
			name:           "fast option combined with enable and enable-all",
			args:           []string{"--enable-all", "--fast", "--enable=unused"},
			enabledLinters: getAllFastLintersWith(t, "unused"),
			noImplicitFast: true,
		},
	}

	testshared.InstallGolangciLint(t)

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			args := []string{"--verbose"}
			if !c.noImplicitFast {
				args = append(args, "--fast")
			}

			r := testshared.NewRunnerBuilder(t).
				WithCommand("linters").
				WithArgs(args...).
				WithArgs(c.args...).
				WithConfig(c.cfg).
				Runner().
				Run()

			sort.Strings(c.enabledLinters)

			r.ExpectOutputContains(fmt.Sprintf("Active %d linters: [%s]",
				len(c.enabledLinters), strings.Join(c.enabledLinters, " ")))
		})
	}
}

func getEnabledByDefaultFastLintersExcept(t *testing.T, except ...string) []string {
	t.Helper()

	m, err := lintersdb.NewManager(nil, nil, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	ebdl := m.GetAllEnabledByDefaultLinters()
	var ret []string
	for _, lc := range ebdl {
		if lc.IsSlowLinter() {
			continue
		}

		if !slices.Contains(except, lc.Name()) {
			ret = append(ret, lc.Name())
		}
	}

	return ret
}

func getAllFastLintersWith(t *testing.T, with ...string) []string {
	t.Helper()

	dbManager, err := lintersdb.NewManager(nil, nil, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	linters := dbManager.GetAllSupportedLinterConfigs()
	ret := append([]string{}, with...)
	for _, lc := range linters {
		if lc.IsSlowLinter() {
			continue
		}
		ret = append(ret, lc.Name())
	}

	return ret
}

func getEnabledByDefaultLinters(t *testing.T) []string {
	t.Helper()

	dbManager, err := lintersdb.NewManager(nil, nil, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	ebdl := dbManager.GetAllEnabledByDefaultLinters()
	var ret []string
	for _, lc := range ebdl {
		if lc.Internal {
			continue
		}

		ret = append(ret, lc.Name())
	}

	return ret
}

func getEnabledByDefaultFastLintersWith(t *testing.T, with ...string) []string {
	t.Helper()

	dbManager, err := lintersdb.NewManager(nil, nil, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	ebdl := dbManager.GetAllEnabledByDefaultLinters()
	ret := append([]string{}, with...)
	for _, lc := range ebdl {
		if lc.IsSlowLinter() {
			continue
		}

		ret = append(ret, lc.Name())
	}

	return ret
}
