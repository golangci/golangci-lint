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

func NewGoCritic(settings *config.GoCriticSettings, cfg *config.Config) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	wrapper := &goCriticWrapper{
		cfg:   cfg,
		sizes: types.SizesFor("gc", runtime.GOARCH),
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
		}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

type goCriticWrapper struct {
	settingsWrapper *goCriticSettingsWrapper
	cfg             *config.Config
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
			logger.Fatalf("%s: %v: setting an explicit GOROOT can fix this problem.", goCriticName, err)
		}
	})

	settingsWrapper := newGoCriticSettingsWrapper(settings, logger)

	settingsWrapper.inferEnabledChecks()

	if err := settingsWrapper.validate(); err != nil {
		logger.Fatalf("%s: invalid settings: %s", goCriticName, err)
	}

	w.settingsWrapper = settingsWrapper
}

func (w *goCriticWrapper) run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	if w.settingsWrapper == nil {
		return nil, fmt.Errorf("the settings wrapper is nil")
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
	allParams := w.settingsWrapper.getLowerCasedParams()

	var enabledCheckers []*gocriticlinter.Checker
	for _, info := range gocriticlinter.GetCheckersInfo() {
		if !w.settingsWrapper.isCheckEnabled(info.Name) {
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

func runGocriticOnPackage(linterCtx *gocriticlinter.Context, checks []*gocriticlinter.Checker,
	files []*ast.File) []result.Issue {
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
		return strings.ReplaceAll(rv.String(), "${configDir}", w.cfg.GetConfigDir())
	default:
		return p
	}
}

// TODO(ldez): rewrite and simplify goCriticSettingsWrapper.

type goCriticSettingsWrapper struct {
	*config.GoCriticSettings

	logger logutils.Log

	allCheckers   []*gocriticlinter.CheckerInfo
	allCheckerMap map[string]*gocriticlinter.CheckerInfo

	inferredEnabledChecks map[string]bool
}

func newGoCriticSettingsWrapper(settings *config.GoCriticSettings, logger logutils.Log) *goCriticSettingsWrapper {
	allCheckers := gocriticlinter.GetCheckersInfo()

	allCheckerMap := make(map[string]*gocriticlinter.CheckerInfo)
	for _, checkInfo := range allCheckers {
		allCheckerMap[checkInfo.Name] = checkInfo
	}

	return &goCriticSettingsWrapper{
		GoCriticSettings:      settings,
		logger:                logger,
		allCheckers:           allCheckers,
		allCheckerMap:         allCheckerMap,
		inferredEnabledChecks: map[string]bool{},
	}
}

func (s *goCriticSettingsWrapper) buildTagToCheckersMap() map[string][]string {
	tagToCheckers := map[string][]string{}

	for _, checker := range s.allCheckers {
		for _, tag := range checker.Tags {
			tagToCheckers[tag] = append(tagToCheckers[tag], checker.Name)
		}
	}

	return tagToCheckers
}

func (s *goCriticSettingsWrapper) checkerTagsDebugf() {
	if !isGoCriticDebug {
		return
	}

	tagToCheckers := s.buildTagToCheckersMap()

	allTags := maps.Keys(tagToCheckers)
	sort.Strings(allTags)

	goCriticDebugf("All gocritic existing tags and checks:")
	for _, tag := range allTags {
		debugChecksListf(tagToCheckers[tag], "  tag %q", tag)
	}
}

func (s *goCriticSettingsWrapper) disabledCheckersDebugf() {
	if !isGoCriticDebug {
		return
	}

	var disabledCheckers []string
	for _, checker := range s.allCheckers {
		if s.inferredEnabledChecks[strings.ToLower(checker.Name)] {
			continue
		}

		disabledCheckers = append(disabledCheckers, checker.Name)
	}

	if len(disabledCheckers) == 0 {
		goCriticDebugf("All checks are enabled")
	} else {
		debugChecksListf(disabledCheckers, "Final not used")
	}
}

func (s *goCriticSettingsWrapper) inferEnabledChecks() {
	s.checkerTagsDebugf()

	enabledByDefaultChecks := s.getDefaultEnabledCheckersNames()
	debugChecksListf(enabledByDefaultChecks, "Enabled by default")

	disabledByDefaultChecks := s.getDefaultDisabledCheckersNames()
	debugChecksListf(disabledByDefaultChecks, "Disabled by default")

	enabledChecks := make([]string, 0, len(s.EnabledTags)+len(enabledByDefaultChecks))

	// EnabledTags
	if len(s.EnabledTags) != 0 {
		tagToCheckers := s.buildTagToCheckersMap()
		for _, tag := range s.EnabledTags {
			enabledChecks = append(enabledChecks, tagToCheckers[tag]...)
		}

		debugChecksListf(enabledChecks, "Enabled by config tags %s", sprintStrings(s.EnabledTags))
	}

	if !(len(s.EnabledTags) == 0 && len(s.EnabledChecks) != 0) {
		// don't use default checks only if we have no enabled tags and enable some checks manually
		enabledChecks = append(enabledChecks, enabledByDefaultChecks...)
	}

	// DisabledTags
	if len(s.DisabledTags) != 0 {
		enabledChecks = s.filterByDisableTags(enabledChecks, s.DisabledTags)
	}

	// EnabledChecks
	if len(s.EnabledChecks) != 0 {
		debugChecksListf(s.EnabledChecks, "Enabled by config")

		alreadyEnabledChecksSet := stringsSliceToSet(enabledChecks)
		for _, enabledCheck := range s.EnabledChecks {
			if alreadyEnabledChecksSet[enabledCheck] {
				s.logger.Warnf("%s: no need to enable check %q: it's already enabled", goCriticName, enabledCheck)
				continue
			}
			enabledChecks = append(enabledChecks, enabledCheck)
		}
	}

	// DisabledChecks
	if len(s.DisabledChecks) != 0 {
		debugChecksListf(s.DisabledChecks, "Disabled by config")

		enabledChecksSet := stringsSliceToSet(enabledChecks)
		for _, disabledCheck := range s.DisabledChecks {
			if !enabledChecksSet[disabledCheck] {
				s.logger.Warnf("%s: check %q was explicitly disabled via config. However, as this check "+
					"is disabled by default, there is no need to explicitly disable it via config.", goCriticName, disabledCheck)
				continue
			}
			delete(enabledChecksSet, disabledCheck)
		}

		enabledChecks = nil
		for enabledCheck := range enabledChecksSet {
			enabledChecks = append(enabledChecks, enabledCheck)
		}
	}

	s.inferredEnabledChecks = map[string]bool{}
	for _, check := range enabledChecks {
		s.inferredEnabledChecks[strings.ToLower(check)] = true
	}

	debugChecksListf(enabledChecks, "Final used")

	s.disabledCheckersDebugf()
}

func (s *goCriticSettingsWrapper) validate() error {
	if len(s.EnabledTags) == 0 {
		if len(s.EnabledChecks) != 0 && len(s.DisabledChecks) != 0 {
			return errors.New("both enabled and disabled check aren't allowed for gocritic")
		}
	} else {
		if err := validateStringsUniq(s.EnabledTags); err != nil {
			return fmt.Errorf("validate enabled tags: %w", err)
		}

		tagToCheckers := s.buildTagToCheckersMap()

		for _, tag := range s.EnabledTags {
			if _, ok := tagToCheckers[tag]; !ok {
				return fmt.Errorf("gocritic [enabled]tag %q doesn't exist", tag)
			}
		}
	}

	if len(s.DisabledTags) > 0 {
		tagToCheckers := s.buildTagToCheckersMap()
		for _, tag := range s.EnabledTags {
			if _, ok := tagToCheckers[tag]; !ok {
				return fmt.Errorf("gocritic [disabled]tag %q doesn't exist", tag)
			}
		}
	}

	if err := validateStringsUniq(s.EnabledChecks); err != nil {
		return fmt.Errorf("validate enabled checks: %w", err)
	}

	if err := validateStringsUniq(s.DisabledChecks); err != nil {
		return fmt.Errorf("validate disabled checks: %w", err)
	}

	if err := s.validateCheckerNames(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (s *goCriticSettingsWrapper) isCheckEnabled(name string) bool {
	return s.inferredEnabledChecks[strings.ToLower(name)]
}

// getAllCheckerNames returns a map containing all checker names supported by gocritic.
func (s *goCriticSettingsWrapper) getAllCheckerNames() map[string]bool {
	allCheckerNames := make(map[string]bool, len(s.allCheckers))

	for _, checker := range s.allCheckers {
		allCheckerNames[strings.ToLower(checker.Name)] = true
	}

	return allCheckerNames
}

func (s *goCriticSettingsWrapper) getDefaultEnabledCheckersNames() []string {
	var enabled []string

	for _, info := range s.allCheckers {
		enable := s.isEnabledByDefaultCheck(info)
		if enable {
			enabled = append(enabled, info.Name)
		}
	}

	return enabled
}

func (s *goCriticSettingsWrapper) getDefaultDisabledCheckersNames() []string {
	var disabled []string

	for _, info := range s.allCheckers {
		enable := s.isEnabledByDefaultCheck(info)
		if !enable {
			disabled = append(disabled, info.Name)
		}
	}

	return disabled
}

func (s *goCriticSettingsWrapper) validateCheckerNames() error {
	allowedNames := s.getAllCheckerNames()

	for _, name := range s.EnabledChecks {
		if !allowedNames[strings.ToLower(name)] {
			return fmt.Errorf("enabled checker %s doesn't exist, all existing checkers: %s",
				name, sprintAllowedCheckerNames(allowedNames))
		}
	}

	for _, name := range s.DisabledChecks {
		if !allowedNames[strings.ToLower(name)] {
			return fmt.Errorf("disabled checker %s doesn't exist, all existing checkers: %s",
				name, sprintAllowedCheckerNames(allowedNames))
		}
	}

	for checkName := range s.SettingsPerCheck {
		if _, ok := allowedNames[checkName]; !ok {
			return fmt.Errorf("invalid setting, checker %s doesn't exist, all existing checkers: %s",
				checkName, sprintAllowedCheckerNames(allowedNames))
		}

		if !s.isCheckEnabled(checkName) {
			s.logger.Warnf("%s: settings were provided for not enabled check %q", goCriticName, checkName)
		}
	}

	return nil
}

func (s *goCriticSettingsWrapper) getLowerCasedParams() map[string]config.GoCriticCheckSettings {
	ret := make(map[string]config.GoCriticCheckSettings, len(s.SettingsPerCheck))

	for checker, params := range s.SettingsPerCheck {
		ret[strings.ToLower(checker)] = params
	}

	return ret
}

func (s *goCriticSettingsWrapper) filterByDisableTags(enabledChecks, disableTags []string) []string {
	enabledChecksSet := stringsSliceToSet(enabledChecks)

	for _, enabledCheck := range enabledChecks {
		checkInfo, checkInfoExists := s.allCheckerMap[enabledCheck]
		if !checkInfoExists {
			s.logger.Warnf("%s: check %q was not exists via filtering disabled tags", goCriticName, enabledCheck)
			continue
		}

		hitTags := intersectStringSlice(checkInfo.Tags, disableTags)
		if len(hitTags) != 0 {
			delete(enabledChecksSet, enabledCheck)
		}
	}

	debugChecksListf(enabledChecks, "Disabled by config tags %s", sprintStrings(disableTags))

	enabledChecks = nil
	for enabledCheck := range enabledChecksSet {
		enabledChecks = append(enabledChecks, enabledCheck)
	}

	return enabledChecks
}

func (s *goCriticSettingsWrapper) isEnabledByDefaultCheck(info *gocriticlinter.CheckerInfo) bool {
	return !info.HasTag("experimental") &&
		!info.HasTag("opinionated") &&
		!info.HasTag("performance")
}

func validateStringsUniq(ss []string) error {
	set := map[string]bool{}

	for _, s := range ss {
		_, ok := set[s]
		if ok {
			return fmt.Errorf("%q occurs multiple times in list", s)
		}
		set[s] = true
	}

	return nil
}

func intersectStringSlice(s1, s2 []string) []string {
	s1Map := make(map[string]struct{}, len(s1))

	for _, s := range s1 {
		s1Map[s] = struct{}{}
	}

	results := make([]string, 0)
	for _, s := range s2 {
		if _, exists := s1Map[s]; exists {
			results = append(results, s)
		}
	}

	return results
}

func sprintAllowedCheckerNames(allowedNames map[string]bool) string {
	namesSlice := maps.Keys(allowedNames)
	return sprintStrings(namesSlice)
}

func sprintStrings(ss []string) string {
	sort.Strings(ss)
	return fmt.Sprint(ss)
}

func debugChecksListf(checks []string, format string, args ...any) {
	if !isGoCriticDebug {
		return
	}

	goCriticDebugf("%s checks (%d): %s", fmt.Sprintf(format, args...), len(checks), sprintStrings(checks))
}

func stringsSliceToSet(ss []string) map[string]bool {
	ret := make(map[string]bool, len(ss))
	for _, s := range ss {
		ret[s] = true
	}

	return ret
}
