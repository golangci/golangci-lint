package golinters

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
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const goCriticName = "gocritic"

var (
	goCriticDebugf  = logutils.Debug(logutils.DebugKeyGoCritic)
	isGoCriticDebug = logutils.HaveDebugTag(logutils.DebugKeyGoCritic)
)

func NewGoCritic(settings *config.GoCriticSettings, lintConfigDir string) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	wrapper := &goCriticWrapper{
		lintConfigDir: lintConfigDir,
		sizes:         types.SizesFor("gc", runtime.GOARCH),
	}

	analyzer := &analysis.Analyzer{
		Name: goCriticName,
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
		goCriticName,
		`Provides diagnostics that check for bugs, performance and style issues.
Extensible without recompilation through dynamic rules.
Dynamic rules are written declaratively with AST patterns, filters, report message and optional suggestion.`,
		[]*analysis.Analyzer{analyzer},
		nil,
	).
		WithContextSetter(func(context *linter.Context) {
			wrapper.init(settings, context.Log)
		}).
		WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
			return resIssues
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

type goCriticWrapper struct {
	settingsWrapper *goCriticSettingsWrapper
	lintConfigDir   string
	sizes           types.Sizes
	once            sync.Once
}

func (w *goCriticWrapper) init(settings *config.GoCriticSettings, logger logutils.Log) {
	if settings == nil {
		return
	}

	w.once.Do(func() {
		err := checkers.InitEmbeddedRules()
		if err != nil {
			logger.Fatalf("%s: %v: setting an explicit GOROOT can fix this problem", goCriticName, err)
		}
	})

	settingsWrapper := newGoCriticSettingsWrapper(settings, logger)
	settingsWrapper.InferEnabledChecks()
	// NOTE(a.telyshev): Validate must be after InferEnabledChecks, not before.
	// Because it uses gathered information about tags set and finally enabled checks.
	if err := settingsWrapper.Validate(); err != nil {
		logger.Fatalf("%s: invalid settings: %s", goCriticName, err)
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

	pkgIssues := runGocriticOnPackage(linterCtx, enabledCheckers, pass.Files)

	issues := make([]goanalysis.Issue, 0, len(pkgIssues))
	for i := range pkgIssues {
		issues = append(issues, goanalysis.NewIssue(&pkgIssues[i], pass))
	}

	return issues, nil
}

func (w *goCriticWrapper) buildEnabledCheckers(linterCtx *gocriticlinter.Context) ([]*gocriticlinter.Checker, error) {
	allParams := w.settingsWrapper.GetLowerCasedParams()

	var enabledCheckers []*gocriticlinter.Checker
	for _, info := range gocriticlinter.GetCheckersInfo() {
		if !w.settingsWrapper.IsCheckEnabled(info.Name) {
			continue
		}

		if err := w.configureCheckerInfo(info, allParams); err != nil {
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

func runGocriticOnPackage(
	linterCtx *gocriticlinter.Context,
	checks []*gocriticlinter.Checker,
	files []*ast.File,
) []result.Issue {
	var res []result.Issue
	for _, f := range files {
		filename := filepath.Base(linterCtx.FileSet.Position(f.Pos()).Filename)
		linterCtx.SetFileInfo(filename, f)

		issues := runGocriticOnFile(linterCtx, f, checks)
		res = append(res, issues...)
	}
	return res
}

func runGocriticOnFile(linterCtx *gocriticlinter.Context, f *ast.File, checks []*gocriticlinter.Checker) []result.Issue {
	var res []result.Issue

	for _, c := range checks {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		for _, warn := range c.Check(f) {
			pos := linterCtx.FileSet.Position(warn.Pos)
			issue := result.Issue{
				Pos:        pos,
				Text:       fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
				FromLinter: goCriticName,
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

func (w *goCriticWrapper) configureCheckerInfo(info *gocriticlinter.CheckerInfo, allParams map[string]config.GoCriticCheckSettings) error {
	params := allParams[strings.ToLower(info.Name)]
	if params == nil { // no config for this checker
		return nil
	}

	infoParams := normalizeCheckerInfoParams(info)
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

func normalizeCheckerInfoParams(info *gocriticlinter.CheckerInfo) gocriticlinter.CheckerParams {
	// lowercase info param keys here because golangci-lint's config parser lowercases all strings
	ret := gocriticlinter.CheckerParams{}
	for k, v := range info.Params {
		ret[strings.ToLower(k)] = v
	}

	return ret
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
		return strings.ReplaceAll(rv.String(), "${configDir}", w.lintConfigDir)
	default:
		return p
	}
}

type goCriticSettingsWrapper struct {
	*config.GoCriticSettings

	logger logutils.Log

	allCheckers       []*gocriticlinter.CheckerInfo
	allCheckersByName map[string]*gocriticlinter.CheckerInfo

	allTagsSorted  []string
	allChecksByTag map[string][]string

	inferredEnabledChecks map[string]struct{}
}

func newGoCriticSettingsWrapper(settings *config.GoCriticSettings, logger logutils.Log) *goCriticSettingsWrapper {
	allCheckers := gocriticlinter.GetCheckersInfo()
	allCheckersByName := make(map[string]*gocriticlinter.CheckerInfo, len(allCheckers))
	for _, checkInfo := range allCheckers {
		allCheckersByName[checkInfo.Name] = checkInfo
	}

	allChecksByTag := make(map[string][]string)
	for _, checker := range allCheckers {
		for _, tag := range checker.Tags {
			allChecksByTag[tag] = append(allChecksByTag[tag], checker.Name)
		}
	}

	allTagsSorted := make([]string, 0, len(allChecksByTag))
	for t := range allChecksByTag {
		allTagsSorted = append(allTagsSorted, t)
	}
	sort.Strings(allTagsSorted)

	return &goCriticSettingsWrapper{
		GoCriticSettings:      settings,
		logger:                logger,
		allCheckers:           allCheckers,
		allCheckersByName:     allCheckersByName,
		allTagsSorted:         allTagsSorted,
		allChecksByTag:        allChecksByTag,
		inferredEnabledChecks: make(map[string]struct{}),
	}
}

func (s *goCriticSettingsWrapper) GetLowerCasedParams() map[string]config.GoCriticCheckSettings {
	ret := make(map[string]config.GoCriticCheckSettings, len(s.SettingsPerCheck))

	for checker, params := range s.SettingsPerCheck {
		ret[strings.ToLower(checker)] = params
	}

	return ret
}

// InferEnabledChecks tries to be consistent with (lintersdb.EnabledSet).build.
func (s *goCriticSettingsWrapper) InferEnabledChecks() {
	s.debugChecksInitialState()

	enabledByDefaultChecks, disabledByDefaultChecks := s.buildEnabledAndDisabledByDefaultChecks()
	debugChecksListf(enabledByDefaultChecks, "Enabled by default")
	debugChecksListf(disabledByDefaultChecks, "Disabled by default")

	enabledChecks := make(map[string]struct{})

	if s.EnableAll {
		enabledChecks = make(map[string]struct{}, len(s.allCheckers))
		for _, info := range s.allCheckers {
			enabledChecks[info.Name] = struct{}{}
		}
	} else if !s.DisableAll {
		// NOTE(a.telyshev): enable-all/disable-all revokes the default settings.
		enabledChecks = make(map[string]struct{}, len(enabledByDefaultChecks))
		for _, check := range enabledByDefaultChecks {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledTags) != 0 {
		enabledFromTags := s.expandTagsToChecks(s.EnabledTags)
		debugChecksListf(enabledFromTags, "Enabled by config tags %s", sprintStrings(s.EnabledTags))

		for _, check := range enabledFromTags {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledChecks) != 0 {
		debugChecksListf(s.EnabledChecks, "Enabled by config")

		for _, check := range s.EnabledChecks {
			if _, ok := enabledChecks[check]; ok {
				s.logger.Warnf("%s: no need to enable check %q: it's already enabled", goCriticName, check)
				continue
			}
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.DisabledTags) != 0 {
		disabledFromTags := s.expandTagsToChecks(s.DisabledTags)
		debugChecksListf(disabledFromTags, "Disabled by config tags %s", sprintStrings(s.DisabledTags))

		for _, check := range disabledFromTags {
			delete(enabledChecks, check)
		}
	}

	if len(s.DisabledChecks) != 0 {
		debugChecksListf(s.DisabledChecks, "Disabled by config")

		for _, check := range s.DisabledChecks {
			if _, ok := enabledChecks[check]; !ok {
				s.logger.Warnf("%s: no need to disable check %q: it's already disabled", goCriticName, check)
				continue
			}
			delete(enabledChecks, check)
		}
	}

	s.inferredEnabledChecks = enabledChecks
	s.debugChecksFinalState()
}

func (s *goCriticSettingsWrapper) buildEnabledAndDisabledByDefaultChecks() (enabled []string, disabled []string) {
	for _, info := range s.allCheckers {
		if enable := isEnabledByDefaultGoCriticChecker(info); enable {
			enabled = append(enabled, info.Name)
		} else {
			disabled = append(disabled, info.Name)
		}
	}
	return enabled, disabled
}

func isEnabledByDefaultGoCriticChecker(info *gocriticlinter.CheckerInfo) bool {
	// https://github.com/go-critic/go-critic/blob/5b67cfd487ae9fe058b4b19321901b3131810f65/cmd/gocritic/check.go#L342-L345
	return !info.HasTag(gocriticlinter.ExperimentalTag) &&
		!info.HasTag(gocriticlinter.OpinionatedTag) &&
		!info.HasTag(gocriticlinter.PerformanceTag) &&
		!info.HasTag(gocriticlinter.SecurityTag)
}

func (s *goCriticSettingsWrapper) expandTagsToChecks(tags []string) []string {
	var checks []string
	for _, t := range tags {
		checks = append(checks, s.allChecksByTag[t]...)
	}
	return checks
}

func (s *goCriticSettingsWrapper) debugChecksInitialState() {
	if !isGoCriticDebug {
		return
	}

	goCriticDebugf("All gocritic existing tags and checks:")
	for _, tag := range s.allTagsSorted {
		debugChecksListf(s.allChecksByTag[tag], "  tag %q", tag)
	}
}

func (s *goCriticSettingsWrapper) debugChecksFinalState() {
	if !isGoCriticDebug {
		return
	}

	var enabledChecks []string
	var disabledChecks []string

	for _, checker := range s.allCheckers {
		name := checker.Name
		if _, ok := s.inferredEnabledChecks[name]; ok {
			enabledChecks = append(enabledChecks, name)
		} else {
			disabledChecks = append(disabledChecks, name)
		}
	}

	debugChecksListf(enabledChecks, "Final used")

	if len(disabledChecks) == 0 {
		goCriticDebugf("All checks are enabled")
	} else {
		debugChecksListf(disabledChecks, "Final not used")
	}
}

func debugChecksListf(checks []string, format string, args ...any) {
	if !isGoCriticDebug {
		return
	}

	goCriticDebugf("%s checks (%d): %s", fmt.Sprintf(format, args...), len(checks), sprintStrings(checks))
}

func sprintStrings(ss []string) string {
	sort.Strings(ss)
	return fmt.Sprint(ss)
}

// Validate tries to be consistent with (lintersdb.Validator).validateEnabledDisabledLintersConfig.
func (s *goCriticSettingsWrapper) Validate() error {
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

func (s *goCriticSettingsWrapper) validateOptionsCombinations() error {
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

func (s *goCriticSettingsWrapper) validateCheckerTags() error {
	for _, tag := range s.EnabledTags {
		if !s.isKnownTag(tag) {
			return fmt.Errorf("enabled tag %q doesn't exist, see %s's documentation", tag, goCriticName)
		}
	}

	for _, tag := range s.DisabledTags {
		if !s.isKnownTag(tag) {
			return fmt.Errorf("disabled tag %q doesn't exist, see %s's documentation", tag, goCriticName)
		}
	}

	return nil
}

func (s *goCriticSettingsWrapper) isKnownTag(tag string) bool {
	_, ok := s.allChecksByTag[tag]
	return ok
}

func (s *goCriticSettingsWrapper) validateCheckerNames() error {
	for _, name := range s.EnabledChecks {
		if !s.isKnownCheck(name) {
			return fmt.Errorf("enabled check %q doesn't exist, see %s's documentation", name, goCriticName)
		}
	}

	for _, name := range s.DisabledChecks {
		if !s.isKnownCheck(name) {
			return fmt.Errorf("disabled check %q doesn't exist, see %s documentation", name, goCriticName)
		}
	}

	for name := range s.SettingsPerCheck {
		if !s.isKnownCheck(name) {
			return fmt.Errorf("invalid settings, check %q doesn't exist, see %s documentation", name, goCriticName)
		}
		if !s.IsCheckEnabled(name) {
			s.logger.Warnf("%s: settings were provided for disabled check %q", goCriticName, name)
		}
	}

	return nil
}

func (s *goCriticSettingsWrapper) isKnownCheck(name string) bool {
	_, ok := s.allCheckersByName[name]
	return ok
}

func (s *goCriticSettingsWrapper) validateDisabledAndEnabledAtOneMoment() error {
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

func (s *goCriticSettingsWrapper) validateAtLeastOneCheckerEnabled() error {
	if len(s.inferredEnabledChecks) == 0 {
		return errors.New("eventually all checks were disabled: at least one must be enabled")
	}
	return nil
}

func (s *goCriticSettingsWrapper) IsCheckEnabled(name string) bool {
	_, ok := s.inferredEnabledChecks[name]
	return ok
}
