package gocritic

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"maps"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/go-critic/go-critic/checkers"
	gocriticlinter "github.com/go-critic/go-critic/linter"
	_ "github.com/quasilyte/go-ruleguard/dsl"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

const linterName = "gocritic"

var (
	debugf  = logutils.Debug(logutils.DebugKeyGoCritic)
	isDebug = logutils.HaveDebugTag(logutils.DebugKeyGoCritic)
)

func New(settings *config.GoCriticSettings) *goanalysis.Linter {
	wrapper := &goCriticWrapper{
		sizes: types.SizesFor("gc", runtime.GOARCH),
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			err := wrapper.run(pass)
			if err != nil {
				return nil, err
			}

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
			wrapper.replacer = strings.NewReplacer(
				internal.PlaceholderBasePath, context.Cfg.GetBasePath(),
			)

			wrapper.init(context.Log, settings)
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

type goCriticWrapper struct {
	settingsWrapper *settingsWrapper
	replacer        *strings.Replacer
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

func (w *goCriticWrapper) run(pass *analysis.Pass) error {
	if w.settingsWrapper == nil {
		return errors.New("the settings wrapper is nil")
	}

	linterCtx := gocriticlinter.NewContext(pass.Fset, w.sizes)

	linterCtx.SetGoVersion(w.settingsWrapper.Go)

	enabledCheckers, err := w.buildEnabledCheckers(linterCtx)
	if err != nil {
		return err
	}

	linterCtx.SetPackageInfo(pass.TypesInfo, pass.Pkg)

	for _, f := range pass.Files {
		runOnFile(pass, f, enabledCheckers)
	}

	return nil
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

		supportedKeys := slices.Sorted(maps.Keys(info.Params))

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
		return w.replacer.Replace(rv.String())
	default:
		return p
	}
}

func runOnFile(pass *analysis.Pass, f *ast.File, checks []*gocriticlinter.Checker) {
	for _, c := range checks {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		for _, warn := range c.Check(f) {
			diag := analysis.Diagnostic{
				Pos:      warn.Pos,
				Category: c.Info.Name,
				Message:  fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
			}

			if warn.HasQuickFix() {
				diag.SuggestedFixes = []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos:     warn.Suggestion.From,
						End:     warn.Suggestion.To,
						NewText: warn.Suggestion.Replacement,
					}},
				}}
			}

			pass.Report(diag)
		}
	}
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
