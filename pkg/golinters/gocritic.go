package golinters

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
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

func NewGocritic() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	sizes := types.SizesFor("gc", runtime.GOARCH)

	analyzer := &analysis.Analyzer{
		Name: gocriticName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gocriticName,
		"The most opinionated Go source code linter",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			linterCtx := gocriticlinter.NewContext(pass.Fset, sizes)
			enabledCheckers, err := buildEnabledCheckers(lintCtx, linterCtx)
			if err != nil {
				return nil, err
			}

			linterCtx.SetPackageInfo(pass.TypesInfo, pass.Pkg)
			var res []goanalysis.Issue
			pkgIssues := runGocriticOnPackage(linterCtx, enabledCheckers, pass.Files)
			for i := range pkgIssues {
				res = append(res, goanalysis.NewIssue(&pkgIssues[i], pass))
			}
			if len(res) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func normalizeCheckerInfoParams(info *gocriticlinter.CheckerInfo) gocriticlinter.CheckerParams {
	// lowercase info param keys here because golangci-lint's config parser lowercases all strings
	ret := gocriticlinter.CheckerParams{}
	for k, v := range info.Params {
		ret[strings.ToLower(k)] = v
	}

	return ret
}

func configureCheckerInfo(info *gocriticlinter.CheckerInfo, allParams map[string]config.GocriticCheckSettings) error {
	params := allParams[strings.ToLower(info.Name)]
	if params == nil { // no config for this checker
		return nil
	}

	infoParams := normalizeCheckerInfoParams(info)
	for k, p := range params {
		v, ok := infoParams[k]
		if ok {
			v.Value = p
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

func buildEnabledCheckers(lintCtx *linter.Context, linterCtx *gocriticlinter.Context) ([]*gocriticlinter.Checker, error) {
	s := lintCtx.Settings().Gocritic
	allParams := s.GetLowercasedParams()

	var enabledCheckers []*gocriticlinter.Checker
	for _, info := range gocriticlinter.GetCheckersInfo() {
		if !s.IsCheckEnabled(info.Name) {
			continue
		}

		if err := configureCheckerInfo(info, allParams); err != nil {
			return nil, err
		}

		c := gocriticlinter.NewChecker(linterCtx, info)
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

func runGocriticOnFile(ctx *gocriticlinter.Context, f *ast.File, checkers []*gocriticlinter.Checker) []result.Issue {
	var res []result.Issue

	for _, c := range checkers {
		// All checkers are expected to use *lint.Context
		// as read-only structure, so no copying is required.
		for _, warn := range c.Check(f) {
			pos := ctx.FileSet.Position(warn.Node.Pos())
			res = append(res, result.Issue{
				Pos:        pos,
				Text:       fmt.Sprintf("%s: %s", c.Info.Name, warn.Text),
				FromLinter: gocriticName,
			})
		}
	}

	return res
}
