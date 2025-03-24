package test

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/test/testshared"
)

func TestEnabledLinters(t *testing.T) {
	// require to display the message "Active x linters: [x,y]"
	t.Setenv(logutils.EnvTestRun, "1")

	cases := []struct {
		name           string
		cfg            string
		enabledLinters []string
		args           []string
	}{
		{
			name: "disable govet in config",
			cfg: `
version: "2"
linters:
	disable:
		- govet
			`,
			enabledLinters: getEnabledByDefaultLintersExcept(t, "govet"),
		},
		{
			name: "enable revive in config",
			cfg: `
version: "2"
linters:
	enable:
		- revive
			`,
			enabledLinters: getEnabledByDefaultLintersWith(t, "revive"),
		},
		{
			name:           "disable govet in cmd",
			args:           []string{"-Dgovet"},
			enabledLinters: getEnabledByDefaultLintersExcept(t, "govet"),
		},
		{
			name: "enable revive in cmd and enable gofmt in config",
			args: []string{"-Erevive"},
			cfg: `
version: "2"
formatters:
	enable:
		- gofmt
			`,
			enabledLinters: getEnabledByDefaultLintersWith(t, "revive", "gofmt"),
		},
		{
			name: "fast option in config",
			cfg: `
version: "2"
linters:
	default: fast
			`,
			enabledLinters: getAllLintersFromGroupFast(t),
		},
		{
			name:           "fast option in flag",
			args:           []string{"--default=fast"},
			enabledLinters: getAllLintersFromGroupFast(t),
		},
		{
			name: "fast option in command-line has higher priority to enable",
			cfg: `
version: "2"
linters:
	default: none
			`,
			args:           []string{"--default=fast"},
			enabledLinters: getAllLintersFromGroupFast(t),
		},
		{
			name:           "only fast linters with standard group",
			args:           []string{"--fast-only"},
			enabledLinters: []string{"ineffassign"},
		},
		{
			name:           "only fast false",
			args:           []string{"--fast-only=false"},
			enabledLinters: getEnabledByDefaultLintersWith(t),
		},
		{
			name:           "fast option combined with enable and enable-all",
			args:           []string{"--default=all", "--fast-only", "--enable=unused"},
			enabledLinters: getAllLintersFromGroupFast(t),
		},
	}

	binPath := testshared.InstallGolangciLint(t)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			args := []string{"--verbose"}

			r := testshared.NewRunnerBuilder(t).
				WithCommand("linters").
				WithArgs(args...).
				WithArgs(c.args...).
				WithConfig(c.cfg).
				WithBinPath(binPath).
				Runner().
				Run()

			sort.Strings(c.enabledLinters)

			r.ExpectOutputContains(fmt.Sprintf("Active %d linters: [%s]",
				len(c.enabledLinters), strings.Join(c.enabledLinters, " ")))
		})
	}
}

func getEnabledByDefaultLintersExcept(t *testing.T, excludes ...string) []string {
	t.Helper()

	linterNames := getEnabledByDefaultLintersWith(t)

	var ret []string
	for _, lc := range linterNames {
		if slices.Contains(excludes, lc) {
			continue
		}

		ret = append(ret, lc)
	}

	return ret
}

func getEnabledByDefaultLintersWith(t *testing.T, includes ...string) []string {
	t.Helper()

	cfg := &config.Config{
		Linters: config.Linters{
			Default: config.GroupStandard,
		},
	}
	dbManager, err := lintersdb.NewManager(logutils.NewStderrLog("skip"), cfg, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	ebdl, err := dbManager.GetEnabledLintersMap()
	require.NoError(t, err)

	var ret []string
	for _, lc := range ebdl {
		if lc.Internal {
			continue
		}

		ret = append(ret, lc.Name())
	}

	return slices.Concat(ret, includes)
}

func getAllLintersFromGroupFast(t *testing.T) []string {
	t.Helper()

	cfg := &config.Config{
		Linters: config.Linters{
			Default: config.GroupAll,
		},
	}

	dbManager, err := lintersdb.NewManager(logutils.NewStderrLog("skip"), cfg, lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	ebdl, err := dbManager.GetEnabledLintersMap()
	require.NoError(t, err)

	var ret []string
	for _, lc := range ebdl {
		if lc.Internal {
			continue
		}

		if lc.IsSlowLinter() {
			continue
		}

		ret = append(ret, lc.Name())
	}

	return ret
}
