package gocritic

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/go-critic/go-critic/checkers"
	gocriticlinter "github.com/go-critic/go-critic/linter"
	_ "github.com/quasilyte/go-ruleguard/dsl"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const linterName = "gocritic"

var (
	debugf  = logutils.Debug(logutils.DebugKeyGoCritic)
	isDebug = logutils.HaveDebugTag(logutils.DebugKeyGoCritic)
)

func New(settings *config.GoCriticSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	wrapper := &goCriticWrapper{
		sizes: types.SizesFor("gc", runtime.GOARCH),
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := wrapper.run(pass)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		`Provides diagnostics that check for bugs, performance and style issues.
Extensible without recompilation through dynamic rules.
Dynamic rules are written declaratively with AST patterns, filters, report message and optional suggestion.`,
		[]*analysis.Analyzer{analyzer},
		nil,
	).
		WithContextSetter(func(context *linter.Context) {
			wrapper.configDir = context.Cfg.GetConfigDir()

			wrapper.init(context.Log, settings)
		}).
		WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
			return resIssues
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

type goCriticWrapper struct {
	settingsWrapper *settingsWrapper
	configDir       string
	sizes           types.Sizes
	once            sync.Once
}

func (w *goCriticWrapper) init(logger logutils.Log, settings *config.GoCriticSettings) {
	if settings == nil {
		return
	}

	w.once.Do(func() {
		err := checkers.InitEmbeddedRules()
		if err != nil {
			logger.Fatalf("%s: %v: setting an explicit GOROOT can fix this problem", linterName, err)
		}
	})

	settingsWrapper := newSettingsWrapper(settings, logger)
	settingsWrapper.InferEnabledChecks()
	// Validate must be after InferEnabledChecks, not before.
	// Because it uses gathered information about tags set and finally enabled checks.
	if err := settingsWrapper.Validate(); err != nil {
		logger.Fatalf("%s: invalid settings: %s", linterName, err)
	}

	w.settingsWrapper = settingsWrapper
}

func (w *goCriticWrapper) run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	if w.settingsWrapper == nil {
		return nil, errors.New("the settings wrapper is nil")
	}

	linterCtx := gocriticlinter.NewContext(pass.Fset, w.sizes)

	linterCtx.SetGoVersion(w.settingsWrapper.Go)

	enabledCheckers, err := w.buildEnabledCheckers(linterCtx)
	if err != nil {
		return nil, err
	}

	linterCtx.SetPackageInfo(pass.TypesInfo, pass.Pkg)

	pkgIssues := runOnPackage(linterCtx, enabledCheckers, pass.Files)

	issues := make([]goanalysis.Issue, 0, len(pkgIssues))
	for i := range pkgIssues {
		issues = append(issues, goanalysis.NewIssue(&pkgIssues[i], pass))
	}

	return issues, nil
}

func (w *goCriticWrapper) buildEnabledCheckers(linterCtx *gocriticlinter.Context) ([]*gocriticlinter.Checker, error) {
	allLowerCasedParams := w.settingsWrapper.GetLowerCasedParams()

	var enabledCheckers []*gocriticlinter.Checker
	for _, info := range gocriticlinter.GetCheckersInfo() {
		if !w.settingsWrapper.IsCheckEnabled(info.Name) {
			continue
		}

		if err := w.configureCheckerInfo(info, allLowerCasedParams); err != nil {
			return nil, err
		}

		c, err := gocriticlinter.NewChecker(linterCtx, info)
		if err != nil {
			return nil, err
		}
		enabledCheckers = append(enabledCheckers, c)
	}

	return enabledCheckers, nil
}

func (w *goCriticWrapper) configureCheckerInfo(
	info *gocriticlinter.CheckerInfo,
	allLowerCasedParams map[string]config.GoCriticCheckSettings,
) error {
	params := allLowerCasedParams[strings.ToLower(info.Name)]
	if params == nil { // no config for this checker
		return nil
	}

	// To lowercase info param keys here because golangci-lint's config parser lowercases all strings.
	infoParams := normalizeMap(info.Params)
	for k, p := range params {
		v, ok := infoParams[k]
		if ok {
			v.Value = w.normalizeCheckerParamsValue(p)
			continue
		}

		// param `k` isn't supported
		if len(info.Params) == 0 {
			return fmt.Errorf("checker %s config param %s doesn't exist: checker doesn't have params",
				info.Name, k)
		}

		supportedKeys := maps.Keys(info.Params)
		sort.Strings(supportedKeys)

		return fmt.Errorf("checker %s config param %s doesn't exist, all existing: %s",
			info.Name, k, supportedKeys)
	}

	return nil
}

// normalizeCheckerParamsValue normalizes value types.
// go-critic asserts that CheckerParam.Value has some specific types,
// but the file parsers (TOML, YAML, JSON) don't create the same representation for raw type.
// then we have to convert value types into the expected value types.
// Maybe in the future, this kind of conversion will be done in go-critic itself.
func (w *goCriticWrapper) normalizeCheckerParamsValue(p any) any {
	rv := reflect.ValueOf(p)
	switch rv.Type().Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return int(rv.Int())
	case reflect.Bool:
		return rv.Bool()
	case reflect.String:
		// Perform variable substitution.
		return strings.ReplaceAll(rv.String(), "${configDir}", w.configDir)
	default:
		return p
	}
}

func runOnPackage(linterCtx *gocriticlinter.Context, checks []*gocriticlinter.Checker, files []*ast.File) []result.Issue {
	var res []result.Issue
	for _, f := range files {
		filename := filepath.Base(linterCtx.FileSet.Position(f.Pos()).Filename)
		linterCtx.SetFileInfo(filename, f)

		issues := runOnFile(linterCtx, f, checks)
		res = append(res, issues...)
	}
	return res
}

func runOnFile(linterCtx *gocriticlinter.Context, f *ast.File, checks []*gocriticlinter.Checker) []result.Issue {
	var res []result.Issue

	for _, c := range checks {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		for _, warn := range c.Check(f) {
			pos := linterCtx.FileSet.Position(warn.Pos)
			issue := result.Issue{
				Pos:        pos,
				Text:       fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
				FromLinter: linterName,
			}

			if warn.HasQuickFix() {
				issue.Replacement = &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  pos.Column - 1,
						Length:    int(warn.Suggestion.To - warn.Suggestion.From),
						NewString: string(warn.Suggestion.Replacement),
					},
				}
			}

			res = append(res, issue)
		}
	}

	return res
}

type goCriticChecks[T any] map[string]T

func (m goCriticChecks[T]) has(name string) bool {
	_, ok := m[name]
	return ok
}

type settingsWrapper struct {
	*config.GoCriticSettings

	logger logutils.Log

	allCheckers []*gocriticlinter.CheckerInfo

	allChecks             goCriticChecks[struct{}]
	allChecksByTag        goCriticChecks[[]string]
	allTagsSorted         []string
	inferredEnabledChecks goCriticChecks[struct{}]

	// *LowerCased fields are used for GoCriticSettings.SettingsPerCheck validation only.

	allChecksLowerCased             goCriticChecks[struct{}]
	inferredEnabledChecksLowerCased goCriticChecks[struct{}]
}

func newSettingsWrapper(settings *config.GoCriticSettings, logger logutils.Log) *settingsWrapper {
	allCheckers := gocriticlinter.GetCheckersInfo()

	allChecks := make(goCriticChecks[struct{}], len(allCheckers))
	allChecksLowerCased := make(goCriticChecks[struct{}], len(allCheckers))
	allChecksByTag := make(goCriticChecks[[]string])
	for _, checker := range allCheckers {
		allChecks[checker.Name] = struct{}{}
		allChecksLowerCased[strings.ToLower(checker.Name)] = struct{}{}

		for _, tag := range checker.Tags {
			allChecksByTag[tag] = append(allChecksByTag[tag], checker.Name)
		}
	}

	allTagsSorted := maps.Keys(allChecksByTag)
	sort.Strings(allTagsSorted)

	return &settingsWrapper{
		GoCriticSettings:                settings,
		logger:                          logger,
		allCheckers:                     allCheckers,
		allChecks:                       allChecks,
		allChecksLowerCased:             allChecksLowerCased,
		allChecksByTag:                  allChecksByTag,
		allTagsSorted:                   allTagsSorted,
		inferredEnabledChecks:           make(goCriticChecks[struct{}]),
		inferredEnabledChecksLowerCased: make(goCriticChecks[struct{}]),
	}
}

func (s *settingsWrapper) IsCheckEnabled(name string) bool {
	return s.inferredEnabledChecks.has(name)
}

func (s *settingsWrapper) GetLowerCasedParams() map[string]config.GoCriticCheckSettings {
	return normalizeMap(s.SettingsPerCheck)
}

// InferEnabledChecks tries to be consistent with (lintersdb.Manager).build.
func (s *settingsWrapper) InferEnabledChecks() {
	s.debugChecksInitialState()

	enabledByDefaultChecks, disabledByDefaultChecks := s.buildEnabledAndDisabledByDefaultChecks()
	debugChecksListf(enabledByDefaultChecks, "Enabled by default")
	debugChecksListf(disabledByDefaultChecks, "Disabled by default")

	enabledChecks := make(goCriticChecks[struct{}])

	if s.EnableAll {
		enabledChecks = make(goCriticChecks[struct{}], len(s.allCheckers))
		for _, info := range s.allCheckers {
			enabledChecks[info.Name] = struct{}{}
		}
	} else if !s.DisableAll {
		// enable-all/disable-all revokes the default settings.
		enabledChecks = make(goCriticChecks[struct{}], len(enabledByDefaultChecks))
		for _, check := range enabledByDefaultChecks {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledTags) != 0 {
		enabledFromTags := s.expandTagsToChecks(s.EnabledTags)
		debugChecksListf(enabledFromTags, "Enabled by config tags %s", sprintSortedStrings(s.EnabledTags))

		for _, check := range enabledFromTags {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledChecks) != 0 {
		debugChecksListf(s.EnabledChecks, "Enabled by config")

		for _, check := range s.EnabledChecks {
			if enabledChecks.has(check) {
				s.logger.Warnf("%s: no need to enable check %q: it's already enabled", linterName, check)
				continue
			}
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.DisabledTags) != 0 {
		disabledFromTags := s.expandTagsToChecks(s.DisabledTags)
		debugChecksListf(disabledFromTags, "Disabled by config tags %s", sprintSortedStrings(s.DisabledTags))

		for _, check := range disabledFromTags {
			delete(enabledChecks, check)
		}
	}

	if len(s.DisabledChecks) != 0 {
		debugChecksListf(s.DisabledChecks, "Disabled by config")

		for _, check := range s.DisabledChecks {
			if !enabledChecks.has(check) {
				s.logger.Warnf("%s: no need to disable check %q: it's already disabled", linterName, check)
				continue
			}
			delete(enabledChecks, check)
		}
	}

	s.inferredEnabledChecks = enabledChecks
	s.inferredEnabledChecksLowerCased = normalizeMap(s.inferredEnabledChecks)
	s.debugChecksFinalState()
}

func (s *settingsWrapper) buildEnabledAndDisabledByDefaultChecks() (enabled, disabled []string) {
	for _, info := range s.allCheckers {
		if enabledByDef := isEnabledByDefaultGoCriticChecker(info); enabledByDef {
			enabled = append(enabled, info.Name)
		} else {
			disabled = append(disabled, info.Name)
		}
	}
	return enabled, disabled
}

func (s *settingsWrapper) expandTagsToChecks(tags []string) []string {
	var checks []string
	for _, tag := range tags {
		checks = append(checks, s.allChecksByTag[tag]...)
	}
	return checks
}

func (s *settingsWrapper) debugChecksInitialState() {
	if !isDebug {
		return
	}

	debugf("All gocritic existing tags and checks:")
	for _, tag := range s.allTagsSorted {
		debugChecksListf(s.allChecksByTag[tag], "  tag %q", tag)
	}
}

func (s *settingsWrapper) debugChecksFinalState() {
	if !isDebug {
		return
	}

	var enabledChecks []string
	var disabledChecks []string

	for _, checker := range s.allCheckers {
		check := checker.Name
		if s.inferredEnabledChecks.has(check) {
			enabledChecks = append(enabledChecks, check)
		} else {
			disabledChecks = append(disabledChecks, check)
		}
	}

	debugChecksListf(enabledChecks, "Final used")

	if len(disabledChecks) == 0 {
		debugf("All checks are enabled")
	} else {
		debugChecksListf(disabledChecks, "Final not used")
	}
}

// Validate tries to be consistent with (lintersdb.Validator).validateEnabledDisabledLintersConfig.
func (s *settingsWrapper) Validate() error {
	for _, v := range []func() error{
		s.validateOptionsCombinations,
		s.validateCheckerTags,
		s.validateCheckerNames,
		s.validateDisabledAndEnabledAtOneMoment,
		s.validateAtLeastOneCheckerEnabled,
	} {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (s *settingsWrapper) validateOptionsCombinations() error {
	if s.EnableAll {
		if s.DisableAll {
			return errors.New("enable-all and disable-all options must not be combined")
		}

		if len(s.EnabledTags) != 0 {
			return errors.New("enable-all and enabled-tags options must not be combined")
		}

		if len(s.EnabledChecks) != 0 {
			return errors.New("enable-all and enabled-checks options must not be combined")
		}
	}

	if s.DisableAll {
		if len(s.DisabledTags) != 0 {
			return errors.New("disable-all and disabled-tags options must not be combined")
		}

		if len(s.DisabledChecks) != 0 {
			return errors.New("disable-all and disabled-checks options must not be combined")
		}

		if len(s.EnabledTags) == 0 && len(s.EnabledChecks) == 0 {
			return errors.New("all checks were disabled, but no one check was enabled: at least one must be enabled")
		}
	}

	return nil
}

func (s *settingsWrapper) validateCheckerTags() error {
	for _, tag := range s.EnabledTags {
		if !s.allChecksByTag.has(tag) {
			return fmt.Errorf("enabled tag %q doesn't exist, see %s's documentation", tag, linterName)
		}
	}

	for _, tag := range s.DisabledTags {
		if !s.allChecksByTag.has(tag) {
			return fmt.Errorf("disabled tag %q doesn't exist, see %s's documentation", tag, linterName)
		}
	}

	return nil
}

func (s *settingsWrapper) validateCheckerNames() error {
	for _, check := range s.EnabledChecks {
		if !s.allChecks.has(check) {
			return fmt.Errorf("enabled check %q doesn't exist, see %s's documentation", check, linterName)
		}
	}

	for _, check := range s.DisabledChecks {
		if !s.allChecks.has(check) {
			return fmt.Errorf("disabled check %q doesn't exist, see %s documentation", check, linterName)
		}
	}

	for check := range s.SettingsPerCheck {
		lcName := strings.ToLower(check)
		if !s.allChecksLowerCased.has(lcName) {
			return fmt.Errorf("invalid check settings: check %q doesn't exist, see %s documentation", check, linterName)
		}
		if !s.inferredEnabledChecksLowerCased.has(lcName) {
			s.logger.Warnf("%s: settings were provided for disabled check %q", check, linterName)
		}
	}

	return nil
}

func (s *settingsWrapper) validateDisabledAndEnabledAtOneMoment() error {
	for _, tag := range s.DisabledTags {
		if slices.Contains(s.EnabledTags, tag) {
			return fmt.Errorf("tag %q disabled and enabled at one moment", tag)
		}
	}

	for _, check := range s.DisabledChecks {
		if slices.Contains(s.EnabledChecks, check) {
			return fmt.Errorf("check %q disabled and enabled at one moment", check)
		}
	}

	return nil
}

func (s *settingsWrapper) validateAtLeastOneCheckerEnabled() error {
	if len(s.inferredEnabledChecks) == 0 {
		return errors.New("eventually all checks were disabled: at least one must be enabled")
	}
	return nil
}

func normalizeMap[ValueT any](in map[string]ValueT) map[string]ValueT {
	ret := make(map[string]ValueT, len(in))
	for k, v := range in {
		ret[strings.ToLower(k)] = v
	}
	return ret
}

func isEnabledByDefaultGoCriticChecker(info *gocriticlinter.CheckerInfo) bool {
	// https://github.com/go-critic/go-critic/blob/5b67cfd487ae9fe058b4b19321901b3131810f65/cmd/gocritic/check.go#L342-L345
	return !info.HasTag(gocriticlinter.ExperimentalTag) &&
		!info.HasTag(gocriticlinter.OpinionatedTag) &&
		!info.HasTag(gocriticlinter.PerformanceTag) &&
		!info.HasTag(gocriticlinter.SecurityTag)
}

func debugChecksListf(checks []string, format string, args ...any) {
	if !isDebug {
		return
	}

	debugf("%s checks (%d): %s", fmt.Sprintf(format, args...), len(checks), sprintSortedStrings(checks))
}

func sprintSortedStrings(v []string) string {
	sort.Strings(slices.Clone(v))
	return fmt.Sprint(v)
}
