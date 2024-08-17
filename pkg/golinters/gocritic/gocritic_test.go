package gocritic

import (
	"strings"
	"testing"

	"github.com/go-critic/go-critic/checkers"
	gocriticlinter "github.com/go-critic/go-critic/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

// https://go-critic.com/overview.html
func Test_settingsWrapper_InferEnabledChecks(t *testing.T) {
	err := checkers.InitEmbeddedRules()
	require.NoError(t, err)

	allCheckersInfo := gocriticlinter.GetCheckersInfo()

	allChecksByTag := make(map[string][]string)
	allChecks := make([]string, 0, len(allCheckersInfo))
	for _, checker := range allCheckersInfo {
		allChecks = append(allChecks, checker.Name)
		for _, tag := range checker.Tags {
			allChecksByTag[tag] = append(allChecksByTag[tag], checker.Name)
		}
	}

	enabledByDefaultChecks := make([]string, 0, len(allCheckersInfo))
	for _, info := range allCheckersInfo {
		if isEnabledByDefaultGoCriticChecker(info) {
			enabledByDefaultChecks = append(enabledByDefaultChecks, info.Name)
		}
	}
	t.Logf("enabled by default checks:\n%s", strings.Join(enabledByDefaultChecks, "\n"))

	insert := func(in []string, toInsert ...string) []string {
		return append(slices.Clone(in), toInsert...)
	}

	remove := func(in []string, toRemove ...string) []string {
		result := slices.Clone(in)
		for _, v := range toRemove {
			if i := slices.Index(result, v); i != -1 {
				result = slices.Delete(result, i, i+1)
			}
		}
		return result
	}

	uniq := func(in []string) []string {
		result := slices.Clone(in)
		slices.Sort(result)
		return slices.Compact(result)
	}

	cases := []struct {
		name                  string
		sett                  *config.GoCriticSettings
		expectedEnabledChecks []string
	}{
		{
			name:                  "no configuration",
			sett:                  &config.GoCriticSettings{},
			expectedEnabledChecks: enabledByDefaultChecks,
		},
		{
			name: "enable checks",
			sett: &config.GoCriticSettings{
				EnabledChecks: []string{"assignOp", "badCall", "emptyDecl"},
			},
			expectedEnabledChecks: insert(enabledByDefaultChecks, "emptyDecl"),
		},
		{
			name: "disable checks",
			sett: &config.GoCriticSettings{
				DisabledChecks: []string{"assignOp", "emptyDecl"},
			},
			expectedEnabledChecks: remove(enabledByDefaultChecks, "assignOp"),
		},
		{
			name: "enable tags",
			sett: &config.GoCriticSettings{
				EnabledTags: []string{"style", "experimental"},
			},
			expectedEnabledChecks: uniq(insert(insert(
				enabledByDefaultChecks,
				allChecksByTag["style"]...),
				allChecksByTag["experimental"]...)),
		},
		{
			name: "disable tags",
			sett: &config.GoCriticSettings{
				DisabledTags: []string{"diagnostic"},
			},
			expectedEnabledChecks: remove(enabledByDefaultChecks, allChecksByTag["diagnostic"]...),
		},
		{
			name: "enable checks disable checks",
			sett: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badLock"},
				DisabledChecks: []string{"assignOp", "badSorting"},
			},
			expectedEnabledChecks: insert(remove(enabledByDefaultChecks, "assignOp"), "badLock"),
		},
		{
			name: "enable checks enable tags",
			sett: &config.GoCriticSettings{
				EnabledChecks: []string{"badCall", "badLock", "hugeParam"},
				EnabledTags:   []string{"diagnostic"},
			},
			expectedEnabledChecks: uniq(insert(insert(enabledByDefaultChecks,
				allChecksByTag["diagnostic"]...),
				"hugeParam")),
		},
		{
			name: "enable checks disable tags",
			sett: &config.GoCriticSettings{
				EnabledChecks: []string{"badCall", "badLock", "boolExprSimplify", "hugeParam"},
				DisabledTags:  []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: insert(remove(remove(enabledByDefaultChecks,
				allChecksByTag["style"]...),
				allChecksByTag["diagnostic"]...),
				"hugeParam"),
		},
		{
			name: "enable all checks via tags",
			sett: &config.GoCriticSettings{
				EnabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedEnabledChecks: allChecks,
		},
		{
			name: "disable checks enable tags",
			sett: &config.GoCriticSettings{
				DisabledChecks: []string{"assignOp", "badCall", "badLock", "hugeParam"},
				EnabledTags:    []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: remove(uniq(insert(insert(enabledByDefaultChecks,
				allChecksByTag["style"]...),
				allChecksByTag["diagnostic"]...)),
				"assignOp", "badCall", "badLock"),
		},
		{
			name: "disable checks disable tags",
			sett: &config.GoCriticSettings{
				DisabledChecks: []string{"badCall", "badLock", "codegenComment", "hugeParam"},
				DisabledTags:   []string{"style"},
			},
			expectedEnabledChecks: remove(remove(enabledByDefaultChecks,
				allChecksByTag["style"]...),
				"badCall", "codegenComment"),
		},
		{
			name: "enable tags disable tags",
			sett: &config.GoCriticSettings{
				EnabledTags:  []string{"experimental"},
				DisabledTags: []string{"style"},
			},
			expectedEnabledChecks: remove(uniq(insert(enabledByDefaultChecks,
				allChecksByTag["experimental"]...)),
				allChecksByTag["style"]...),
		},
		{
			name: "enable checks disable checks enable tags",
			sett: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badLock", "boolExprSimplify", "indexAlloc", "hugeParam"},
				DisabledChecks: []string{"deprecatedComment", "typeSwitchVar"},
				EnabledTags:    []string{"experimental"},
			},
			expectedEnabledChecks: remove(uniq(insert(insert(enabledByDefaultChecks,
				allChecksByTag["experimental"]...),
				"indexAlloc", "hugeParam")),
				"deprecatedComment", "typeSwitchVar"),
		},
		{
			name: "enable checks disable checks enable tags disable tags",
			sett: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badCond", "badLock", "indexAlloc", "hugeParam"},
				DisabledChecks: []string{"deprecatedComment", "typeSwitchVar"},
				EnabledTags:    []string{"experimental"},
				DisabledTags:   []string{"performance"},
			},
			expectedEnabledChecks: remove(remove(uniq(insert(insert(enabledByDefaultChecks,
				allChecksByTag["experimental"]...),
				"badCond")),
				allChecksByTag["performance"]...),
				"deprecatedComment", "typeSwitchVar"),
		},
		{
			name: "enable single tag only",
			sett: &config.GoCriticSettings{
				DisableAll:  true,
				EnabledTags: []string{"experimental"},
			},
			expectedEnabledChecks: allChecksByTag["experimental"],
		},
		{
			name: "enable two tags only",
			sett: &config.GoCriticSettings{
				DisableAll:  true,
				EnabledTags: []string{"experimental", "performance"},
			},
			expectedEnabledChecks: uniq(insert(allChecksByTag["experimental"], allChecksByTag["performance"]...)),
		},
		{
			name: "disable single tag only",
			sett: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"style"},
			},
			expectedEnabledChecks: remove(allChecks, allChecksByTag["style"]...),
		},
		{
			name: "disable two tags only",
			sett: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: remove(remove(allChecks, allChecksByTag["style"]...), allChecksByTag["diagnostic"]...),
		},
		{
			name: "enable some checks only",
			sett: &config.GoCriticSettings{
				DisableAll:    true,
				EnabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
			},
			expectedEnabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
		},
		{
			name: "disable some checks only",
			sett: &config.GoCriticSettings{
				EnableAll:      true,
				DisabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
			},
			expectedEnabledChecks: remove(allChecks, "deferInLoop", "dupImport", "ifElseChain", "mapKey"),
		},
		{
			name: "enable single tag and some checks from another tag only",
			sett: &config.GoCriticSettings{
				DisableAll:    true,
				EnabledTags:   []string{"experimental"},
				EnabledChecks: []string{"importShadow"},
			},
			expectedEnabledChecks: insert(allChecksByTag["experimental"], "importShadow"),
		},
		{
			name: "disable single tag and some checks from another tag only",
			sett: &config.GoCriticSettings{
				EnableAll:      true,
				DisabledTags:   []string{"experimental"},
				DisabledChecks: []string{"importShadow"},
			},
			expectedEnabledChecks: remove(remove(allChecks, allChecksByTag["experimental"]...), "importShadow"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lg := logutils.NewStderrLog("Test_goCriticSettingsWrapper_InferEnabledChecks")
			wr := newSettingsWrapper(tt.sett, lg)

			wr.InferEnabledChecks()
			assert.ElementsMatch(t, tt.expectedEnabledChecks, maps.Keys(wr.inferredEnabledChecks))
			assert.NoError(t, wr.Validate())
		})
	}
}

func Test_settingsWrapper_Validate(t *testing.T) {
	cases := []struct {
		name        string
		sett        *config.GoCriticSettings
		expectedErr bool
	}{
		{
			name: "combine enable-all and disable-all",
			sett: &config.GoCriticSettings{
				EnableAll:  true,
				DisableAll: true,
			},
			expectedErr: true,
		},
		{
			name: "combine enable-all and enabled-tags",
			sett: &config.GoCriticSettings{
				EnableAll:   true,
				EnabledTags: []string{"experimental"},
			},
			expectedErr: true,
		},
		{
			name: "combine enable-all and enabled-checks",
			sett: &config.GoCriticSettings{
				EnableAll:     true,
				EnabledChecks: []string{"dupImport"},
			},
			expectedErr: true,
		},
		{
			name: "combine disable-all and disabled-tags",
			sett: &config.GoCriticSettings{
				DisableAll:   true,
				DisabledTags: []string{"style"},
			},
			expectedErr: true,
		},
		{
			name: "combine disable-all and disable-checks",
			sett: &config.GoCriticSettings{
				DisableAll:     true,
				DisabledChecks: []string{"appendAssign"},
			},
			expectedErr: true,
		},
		{
			name: "disable-all and no one check enabled",
			sett: &config.GoCriticSettings{
				DisableAll: true,
			},
			expectedErr: true,
		},
		{
			name: "unknown enabled tag",
			sett: &config.GoCriticSettings{
				EnabledTags: []string{"diagnostic", "go-proverbs"},
			},
			expectedErr: true,
		},
		{
			name: "unknown disabled tag",
			sett: &config.GoCriticSettings{
				DisabledTags: []string{"style", "go-proverbs"},
			},
			expectedErr: true,
		},
		{
			name: "unknown enabled check",
			sett: &config.GoCriticSettings{
				EnabledChecks: []string{"appendAssign", "noExitAfterDefer", "underef"},
			},
			expectedErr: true,
		},
		{
			name: "unknown disabled check",
			sett: &config.GoCriticSettings{
				DisabledChecks: []string{"dupSubExpr", "noExitAfterDefer", "returnAfterHttpError"},
			},
			expectedErr: true,
		},
		{
			name: "settings for unknown check",
			sett: &config.GoCriticSettings{
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"captLocall":    {"paramsOnly": false},
					"unnamedResult": {"checkExported": true},
				},
			},
			expectedErr: true,
		},
		{
			name: "settings for disabled check",
			sett: &config.GoCriticSettings{
				DisabledChecks: []string{"elseif"},
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"elseif": {"skipBalanced": true},
				},
			},
			expectedErr: false, // Just logging.
		},
		{
			name: "settings by lower-cased checker name",
			sett: &config.GoCriticSettings{
				EnabledChecks: []string{"tooManyResultsChecker"},
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"toomanyresultschecker": {"maxResults": 3},
					"unnamedResult":         {"checkExported": true},
				},
			},
			expectedErr: false,
		},
		{
			name: "enabled and disabled at one moment check",
			sett: &config.GoCriticSettings{
				EnabledChecks:  []string{"appendAssign", "codegenComment", "underef"},
				DisabledChecks: []string{"elseif", "underef"},
			},
			expectedErr: true,
		},
		{
			name: "enabled and disabled at one moment tag",
			sett: &config.GoCriticSettings{
				EnabledTags:  []string{"performance", "style"},
				DisabledTags: []string{"style", "diagnostic"},
			},
			expectedErr: true,
		},
		{
			name: "disable all checks via tags",
			sett: &config.GoCriticSettings{
				DisabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedErr: true,
		},
		{
			name: "enable-all and disable all checks via tags",
			sett: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedErr: true,
		},
		{
			name: "valid configuration",
			sett: &config.GoCriticSettings{
				EnabledTags:    []string{"performance"},
				DisabledChecks: []string{"dupImport", "ifElseChain", "octalLiteral", "whyNoLint"},
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"hugeParam":    {"sizeThreshold": 100},
					"rangeValCopy": {"skipTestFuncs": true},
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lg := logutils.NewStderrLog("Test_goCriticSettingsWrapper_Validate")
			wr := newSettingsWrapper(tt.sett, lg)

			wr.InferEnabledChecks()

			err := wr.Validate()
			if tt.expectedErr {
				if assert.Error(t, err) {
					t.Log(err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
