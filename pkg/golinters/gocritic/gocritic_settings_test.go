package gocritic

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/go-critic/go-critic/checkers"
	gocriticlinter "github.com/go-critic/go-critic/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

// https://go-critic.com/overview.html
func Test_settingsWrapper_inferEnabledChecks(t *testing.T) {
	err := checkers.InitEmbeddedRules()
	require.NoError(t, err)

	allCheckersInfo := gocriticlinter.GetCheckersInfo()

	allChecksByTag := make(map[string]Slicer)
	allChecks := make(Slicer, 0, len(allCheckersInfo))

	for _, checker := range allCheckersInfo {
		allChecks = append(allChecks, checker.Name)
		for _, tag := range checker.Tags {
			allChecksByTag[tag] = append(allChecksByTag[tag], checker.Name)
		}
	}

	enabledByDefaultChecks := make(Slicer, 0, len(allCheckersInfo))

	for _, info := range allCheckersInfo {
		if isEnabledByDefaultGoCriticChecker(info) {
			enabledByDefaultChecks = append(enabledByDefaultChecks, info.Name)
		}
	}

	t.Logf("enabled by default checks:\n%s", strings.Join(enabledByDefaultChecks, "\n"))

	testCases := []struct {
		name                  string
		settings              *config.GoCriticSettings
		expectedEnabledChecks []string
	}{
		{
			name:                  "no configuration",
			settings:              &config.GoCriticSettings{},
			expectedEnabledChecks: enabledByDefaultChecks,
		},
		{
			name: "enable checks",
			settings: &config.GoCriticSettings{
				EnabledChecks: []string{"assignOp", "badCall", "emptyDecl"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.add("emptyDecl"),
		},
		{
			name: "disable checks",
			settings: &config.GoCriticSettings{
				DisabledChecks: []string{"assignOp", "emptyDecl"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.remove("assignOp"),
		},
		{
			name: "enable tags",
			settings: &config.GoCriticSettings{
				EnabledTags: []string{"style", "experimental"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["style"]...).
				add(allChecksByTag["experimental"]...).
				uniq(),
		},
		{
			name: "disable tags",
			settings: &config.GoCriticSettings{
				DisabledTags: []string{"diagnostic"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.remove(allChecksByTag["diagnostic"]...),
		},
		{
			name: "enable checks disable checks",
			settings: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badLock"},
				DisabledChecks: []string{"assignOp", "badSorting"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				remove("assignOp").
				add("badLock"),
		},
		{
			name: "enable checks enable tags",
			settings: &config.GoCriticSettings{
				EnabledChecks: []string{"badCall", "badLock", "hugeParam"},
				EnabledTags:   []string{"diagnostic"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["diagnostic"]...).
				add("hugeParam").
				uniq(),
		},
		{
			name: "enable checks disable tags",
			settings: &config.GoCriticSettings{
				EnabledChecks: []string{"badCall", "badLock", "boolExprSimplify", "hugeParam"},
				DisabledTags:  []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				remove(allChecksByTag["style"]...).
				remove(allChecksByTag["diagnostic"]...).
				add("hugeParam"),
		},
		{
			name: "enable all checks via tags",
			settings: &config.GoCriticSettings{
				EnabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedEnabledChecks: allChecks,
		},
		{
			name: "disable checks enable tags",
			settings: &config.GoCriticSettings{
				DisabledChecks: []string{"assignOp", "badCall", "badLock", "hugeParam"},
				EnabledTags:    []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["style"]...).
				add(allChecksByTag["diagnostic"]...).
				uniq().
				remove("assignOp", "badCall", "badLock"),
		},
		{
			name: "disable checks disable tags",
			settings: &config.GoCriticSettings{
				DisabledChecks: []string{"badCall", "badLock", "codegenComment", "hugeParam"},
				DisabledTags:   []string{"style"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				remove(allChecksByTag["style"]...).
				remove("badCall", "codegenComment"),
		},
		{
			name: "enable tags disable tags",
			settings: &config.GoCriticSettings{
				EnabledTags:  []string{"experimental"},
				DisabledTags: []string{"style"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["experimental"]...).
				uniq().
				remove(allChecksByTag["style"]...),
		},
		{
			name: "enable checks disable checks enable tags",
			settings: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badLock", "boolExprSimplify", "indexAlloc", "hugeParam"},
				DisabledChecks: []string{"deprecatedComment", "typeSwitchVar"},
				EnabledTags:    []string{"experimental"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["experimental"]...).
				add("indexAlloc", "hugeParam").
				uniq().
				remove("deprecatedComment", "typeSwitchVar"),
		},
		{
			name: "enable checks disable checks enable tags disable tags",
			settings: &config.GoCriticSettings{
				EnabledChecks:  []string{"badCall", "badCond", "badLock", "indexAlloc", "hugeParam"},
				DisabledChecks: []string{"deprecatedComment", "typeSwitchVar"},
				EnabledTags:    []string{"experimental"},
				DisabledTags:   []string{"performance"},
			},
			expectedEnabledChecks: enabledByDefaultChecks.
				add(allChecksByTag["experimental"]...).
				add("badCond").
				uniq().
				remove(allChecksByTag["performance"]...).
				remove("deprecatedComment", "typeSwitchVar"),
		},
		{
			name: "enable single tag only",
			settings: &config.GoCriticSettings{
				DisableAll:  true,
				EnabledTags: []string{"experimental"},
			},
			expectedEnabledChecks: allChecksByTag["experimental"],
		},
		{
			name: "enable two tags only",
			settings: &config.GoCriticSettings{
				DisableAll:  true,
				EnabledTags: []string{"experimental", "performance"},
			},
			expectedEnabledChecks: allChecksByTag["experimental"].
				add(allChecksByTag["performance"]...).
				uniq(),
		},
		{
			name: "disable single tag only",
			settings: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"style"},
			},
			expectedEnabledChecks: allChecks.remove(allChecksByTag["style"]...),
		},
		{
			name: "disable two tags only",
			settings: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"style", "diagnostic"},
			},
			expectedEnabledChecks: allChecks.
				remove(allChecksByTag["style"]...).
				remove(allChecksByTag["diagnostic"]...),
		},
		{
			name: "enable some checks only",
			settings: &config.GoCriticSettings{
				DisableAll:    true,
				EnabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
			},
			expectedEnabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
		},
		{
			name: "disable some checks only",
			settings: &config.GoCriticSettings{
				EnableAll:      true,
				DisabledChecks: []string{"deferInLoop", "dupImport", "ifElseChain", "mapKey"},
			},
			expectedEnabledChecks: allChecks.
				remove("deferInLoop", "dupImport", "ifElseChain", "mapKey"),
		},
		{
			name: "enable single tag and some checks from another tag only",
			settings: &config.GoCriticSettings{
				DisableAll:    true,
				EnabledTags:   []string{"experimental"},
				EnabledChecks: []string{"importShadow"},
			},
			expectedEnabledChecks: allChecksByTag["experimental"].add("importShadow"),
		},
		{
			name: "disable single tag and some checks from another tag only",
			settings: &config.GoCriticSettings{
				EnableAll:      true,
				DisabledTags:   []string{"experimental"},
				DisabledChecks: []string{"importShadow"},
			},
			expectedEnabledChecks: allChecks.
				remove(allChecksByTag["experimental"]...).
				remove("importShadow"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			lg := logutils.NewStderrLog(t.Name())
			wr := newSettingsWrapper(lg, test.settings, nil)

			wr.inferEnabledChecks()

			assert.ElementsMatch(t, test.expectedEnabledChecks, slices.Collect(maps.Keys(wr.inferredEnabledChecks)))

			assert.NoError(t, wr.validate())
		})
	}
}

func Test_settingsWrapper_Load(t *testing.T) {
	testCases := []struct {
		name        string
		settings    *config.GoCriticSettings
		expectedErr bool
	}{
		{
			name: "combine enable-all and disable-all",
			settings: &config.GoCriticSettings{
				EnableAll:  true,
				DisableAll: true,
			},
			expectedErr: true,
		},
		{
			name: "combine enable-all and enabled-tags",
			settings: &config.GoCriticSettings{
				EnableAll:   true,
				EnabledTags: []string{"experimental"},
			},
			expectedErr: true,
		},
		{
			name: "combine enable-all and enabled-checks",
			settings: &config.GoCriticSettings{
				EnableAll:     true,
				EnabledChecks: []string{"dupImport"},
			},
			expectedErr: true,
		},
		{
			name: "combine disable-all and disabled-tags",
			settings: &config.GoCriticSettings{
				DisableAll:   true,
				DisabledTags: []string{"style"},
			},
			expectedErr: true,
		},
		{
			name: "combine disable-all and disable-checks",
			settings: &config.GoCriticSettings{
				DisableAll:     true,
				DisabledChecks: []string{"appendAssign"},
			},
			expectedErr: true,
		},
		{
			name: "disable-all and no one check enabled",
			settings: &config.GoCriticSettings{
				DisableAll: true,
			},
			expectedErr: true,
		},
		{
			name: "unknown enabled tag",
			settings: &config.GoCriticSettings{
				EnabledTags: []string{"diagnostic", "go-proverbs"},
			},
			expectedErr: true,
		},
		{
			name: "unknown disabled tag",
			settings: &config.GoCriticSettings{
				DisabledTags: []string{"style", "go-proverbs"},
			},
			expectedErr: true,
		},
		{
			name: "unknown enabled check",
			settings: &config.GoCriticSettings{
				EnabledChecks: []string{"appendAssign", "noExitAfterDefer", "underef"},
			},
			expectedErr: true,
		},
		{
			name: "unknown disabled check",
			settings: &config.GoCriticSettings{
				DisabledChecks: []string{"dupSubExpr", "noExitAfterDefer", "returnAfterHttpError"},
			},
			expectedErr: true,
		},
		{
			name: "settings for unknown check",
			settings: &config.GoCriticSettings{
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"captLocall":    {"paramsOnly": false},
					"unnamedResult": {"checkExported": true},
				},
			},
			expectedErr: true,
		},
		{
			name: "settings for disabled check",
			settings: &config.GoCriticSettings{
				DisabledChecks: []string{"elseif"},
				SettingsPerCheck: map[string]config.GoCriticCheckSettings{
					"elseif": {"skipBalanced": true},
				},
			},
			expectedErr: false, // Just logging.
		},
		{
			name: "settings by lower-cased checker name",
			settings: &config.GoCriticSettings{
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
			settings: &config.GoCriticSettings{
				EnabledChecks:  []string{"appendAssign", "codegenComment", "underef"},
				DisabledChecks: []string{"elseif", "underef"},
			},
			expectedErr: true,
		},
		{
			name: "enabled and disabled at one moment tag",
			settings: &config.GoCriticSettings{
				EnabledTags:  []string{"performance", "style"},
				DisabledTags: []string{"style", "diagnostic"},
			},
			expectedErr: true,
		},
		{
			name: "disable all checks via tags",
			settings: &config.GoCriticSettings{
				DisabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedErr: true,
		},
		{
			name: "enable-all and disable all checks via tags",
			settings: &config.GoCriticSettings{
				EnableAll:    true,
				DisabledTags: []string{"diagnostic", "experimental", "opinionated", "performance", "style"},
			},
			expectedErr: true,
		},
		{
			name: "valid configuration",
			settings: &config.GoCriticSettings{
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

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			lg := logutils.NewStderrLog(t.Name())
			wr := newSettingsWrapper(lg, test.settings, nil)

			err := wr.Load()
			if test.expectedErr {
				if assert.Error(t, err) {
					t.Log(err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type Slicer []string

func (s Slicer) add(toAdd ...string) Slicer {
	return slices.Concat(s, toAdd)
}

func (s Slicer) remove(toRemove ...string) Slicer {
	result := slices.Clone(s)

	for _, v := range toRemove {
		if i := slices.Index(result, v); i != -1 {
			result = slices.Delete(result, i, i+1)
		}
	}

	return result
}

func (s Slicer) uniq() Slicer {
	return slices.Compact(slices.Sorted(slices.Values(s)))
}
