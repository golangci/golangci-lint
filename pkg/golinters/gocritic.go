package golinters

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	gocriticlinter "github.com/go-critic/go-critic/framework/linter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const gocriticName = "gocritic"

func NewGocritic(settings *config.GocriticSettings, cfg *config.Config) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	sizes := types.SizesFor("gc", runtime.GOARCH)

	wrapper := goCriticWrapper{
		settings: settings,
		cfg:      cfg,
		sizes:    sizes,
	}

	analyzer := &analysis.Analyzer{
		Name: gocriticName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
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
		gocriticName,
		`Provides diagnostics that check for bugs, performance and style issues.
Extensible without recompilation through dynamic rules.
Dynamic rules are written declaratively with AST patterns, filters, report message and optional suggestion.`,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

type goCriticWrapper struct {
	settings *config.GocriticSettings
	cfg      *config.Config
	sizes    types.Sizes
}

func (w goCriticWrapper) run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	linterCtx := gocriticlinter.NewContext(pass.Fset, w.sizes)

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

func (w goCriticWrapper) buildEnabledCheckers(linterCtx *gocriticlinter.Context) ([]*gocriticlinter.Checker, error) {
	allParams := w.settings.GetLowercasedParams()

	var enabledCheckers []*gocriticlinter.Checker
	for _, info := range gocriticlinter.GetCheckersInfo() {
		if !w.settings.IsCheckEnabled(info.Name) {
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

func runGocriticOnPackage(linterCtx *gocriticlinter.Context, checkers []*gocriticlinter.Checker,
	files []*ast.File) []result.Issue {
	var res []result.Issue
	for _, f := range files {
		filename := filepath.Base(linterCtx.FileSet.Position(f.Pos()).Filename)
		linterCtx.SetFileInfo(filename, f)

		issues := runGocriticOnFile(linterCtx, f, checkers)
		res = append(res, issues...)
	}
	return res
}

func runGocriticOnFile(linterCtx *gocriticlinter.Context, f *ast.File, checkers []*gocriticlinter.Checker) []result.Issue {
	var res []result.Issue

	for _, c := range checkers {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		for _, warn := range c.Check(f) {
			pos := linterCtx.FileSet.Position(warn.Node.Pos())
			issue := result.Issue{
				Pos:        pos,
				Text:       fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
				FromLinter: gocriticName,
			}

			if warn.HasQuickFix() {
				issue.Replacement = &result.Replacement{
					Inline: &result.InlineFix{
						StartCol:  pos.Column - 1,
						Length:    int(warn.Node.End() - warn.Node.Pos()),
						NewString: string(warn.Suggestion.Replacement),
					},
				}
			}

			res = append(res, issue)
		}
	}

	return res
}

func (w goCriticWrapper) configureCheckerInfo(info *gocriticlinter.CheckerInfo, allParams map[string]config.GocriticCheckSettings) error {
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

		var supportedKeys []string
		for sk := range info.Params {
			supportedKeys = append(supportedKeys, sk)
		}
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
func (w goCriticWrapper) normalizeCheckerParamsValue(p interface{}) interface{} {
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
