package goanalysis

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golangci/golangci-lint/pkg/timeutils"

	"github.com/golangci/golangci-lint/internal/pkgcache"
	"github.com/golangci/golangci-lint/pkg/logutils"

	"golang.org/x/tools/go/packages"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	libpackages "github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	TheOnlyAnalyzerName = "the_only_name"
	TheOnlyanalyzerDoc  = "the_only_doc"
)

type LoadMode int

const (
	LoadModeNone LoadMode = iota
	LoadModeSyntax
	LoadModeTypesInfo
	LoadModeWholeProgram
)

var issuesCacheDebugf = logutils.Debug("goanalysis/issues/cache")

func (loadMode LoadMode) String() string {
	switch loadMode {
	case LoadModeNone:
		return "none"
	case LoadModeSyntax:
		return "syntax"
	case LoadModeTypesInfo:
		return "types info"
	case LoadModeWholeProgram:
		return "whole program"
	}
	panic(fmt.Sprintf("unknown load mode %d", loadMode))
}

type Linter struct {
	name, desc              string
	analyzers               []*analysis.Analyzer
	cfg                     map[string]map[string]interface{}
	issuesReporter          func(*linter.Context) []Issue
	contextSetter           func(*linter.Context)
	loadMode                LoadMode
	needUseOriginalPackages bool
	isTypecheckModeOn       bool
}

func NewLinter(name, desc string, analyzers []*analysis.Analyzer, cfg map[string]map[string]interface{}) *Linter {
	return &Linter{name: name, desc: desc, analyzers: analyzers, cfg: cfg}
}

func (lnt *Linter) UseOriginalPackages() {
	lnt.needUseOriginalPackages = true
}

func (lnt *Linter) SetTypecheckMode() {
	lnt.isTypecheckModeOn = true
}

func (lnt *Linter) LoadMode() LoadMode {
	return lnt.loadMode
}

func (lnt *Linter) WithLoadMode(loadMode LoadMode) *Linter {
	lnt.loadMode = loadMode
	return lnt
}

func (lnt *Linter) WithIssuesReporter(r func(*linter.Context) []Issue) *Linter {
	lnt.issuesReporter = r
	return lnt
}

func (lnt *Linter) WithContextSetter(cs func(*linter.Context)) *Linter {
	lnt.contextSetter = cs
	return lnt
}

func (lnt *Linter) Name() string {
	return lnt.name
}

func (lnt *Linter) Desc() string {
	return lnt.desc
}

func (lnt *Linter) allAnalyzerNames() []string {
	var ret []string
	for _, a := range lnt.analyzers {
		ret = append(ret, a.Name)
	}
	return ret
}

func allFlagNames(fs *flag.FlagSet) []string {
	var ret []string
	fs.VisitAll(func(f *flag.Flag) {
		ret = append(ret, f.Name)
	})
	return ret
}

func valueToString(v interface{}) string {
	if ss, ok := v.([]string); ok {
		return strings.Join(ss, ",")
	}

	if is, ok := v.([]interface{}); ok {
		var ss []string
		for _, i := range is {
			ss = append(ss, fmt.Sprint(i))
		}
		return valueToString(ss)
	}

	return fmt.Sprint(v)
}

func (lnt *Linter) configureAnalyzer(a *analysis.Analyzer, cfg map[string]interface{}) error {
	for k, v := range cfg {
		f := a.Flags.Lookup(k)
		if f == nil {
			validFlagNames := allFlagNames(&a.Flags)
			if len(validFlagNames) == 0 {
				return fmt.Errorf("analyzer doesn't have settings")
			}

			return fmt.Errorf("analyzer doesn't have setting %q, valid settings: %v",
				k, validFlagNames)
		}

		if err := f.Value.Set(valueToString(v)); err != nil {
			return errors.Wrapf(err, "failed to set analyzer setting %q with value %v", k, v)
		}
	}

	return nil
}

func (lnt *Linter) configure() error {
	analyzersMap := map[string]*analysis.Analyzer{}
	for _, a := range lnt.analyzers {
		analyzersMap[a.Name] = a
	}

	for analyzerName, analyzerSettings := range lnt.cfg {
		a := analyzersMap[analyzerName]
		if a == nil {
			return fmt.Errorf("settings key %q must be valid analyzer name, valid analyzers: %v",
				analyzerName, lnt.allAnalyzerNames())
		}

		if err := lnt.configureAnalyzer(a, analyzerSettings); err != nil {
			return errors.Wrapf(err, "failed to configure analyzer %s", analyzerName)
		}
	}

	return nil
}

func parseError(srcErr packages.Error) (*result.Issue, error) {
	pos, err := libpackages.ParseErrorPosition(srcErr.Pos)
	if err != nil {
		return nil, err
	}

	return &result.Issue{
		Pos:        *pos,
		Text:       srcErr.Msg,
		FromLinter: "typecheck",
	}, nil
}

func buildIssuesFromErrorsForTypecheckMode(errs []error, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []result.Issue
	uniqReportedIssues := map[string]bool{}
	for _, err := range errs {
		itErr, ok := errors.Cause(err).(*IllTypedError)
		if !ok {
			return nil, err
		}
		for _, err := range libpackages.ExtractErrors(itErr.Pkg) {
			i, perr := parseError(err)
			if perr != nil { // failed to parse
				if uniqReportedIssues[err.Msg] {
					continue
				}
				uniqReportedIssues[err.Msg] = true
				lintCtx.Log.Errorf("typechecking error: %s", err.Msg)
			} else {
				i.Pkg = itErr.Pkg // to save to cache later
				issues = append(issues, *i)
			}
		}
	}
	return issues, nil
}

func buildIssues(diags []Diagnostic, linterNameBuilder func(diag *Diagnostic) string) []result.Issue {
	var issues []result.Issue
	for i := range diags {
		diag := &diags[i]
		linterName := linterNameBuilder(diag)
		var text string
		if diag.Analyzer.Name == linterName {
			text = diag.Message
		} else {
			text = fmt.Sprintf("%s: %s", diag.Analyzer.Name, diag.Message)
		}
		issues = append(issues, result.Issue{
			FromLinter: linterName,
			Text:       text,
			Pos:        diag.Position,
			Pkg:        diag.Pkg,
		})
	}
	return issues
}

func (lnt *Linter) preRun(lintCtx *linter.Context) error {
	if err := analysis.Validate(lnt.analyzers); err != nil {
		return errors.Wrap(err, "failed to validate analyzers")
	}

	if err := lnt.configure(); err != nil {
		return errors.Wrap(err, "failed to configure analyzers")
	}

	if lnt.contextSetter != nil {
		lnt.contextSetter(lintCtx)
	}

	return nil
}

func (lnt *Linter) getName() string {
	return lnt.name
}

func (lnt *Linter) getLinterNameForDiagnostic(*Diagnostic) string {
	return lnt.name
}

func (lnt *Linter) getAnalyzers() []*analysis.Analyzer {
	return lnt.analyzers
}

func (lnt *Linter) useOriginalPackages() bool {
	return lnt.needUseOriginalPackages
}

func (lnt *Linter) isTypecheckMode() bool {
	return lnt.isTypecheckModeOn
}

func (lnt *Linter) reportIssues(lintCtx *linter.Context) []Issue {
	if lnt.issuesReporter != nil {
		return lnt.issuesReporter(lintCtx)
	}
	return nil
}

func (lnt *Linter) getLoadMode() LoadMode {
	return lnt.loadMode
}

type runAnalyzersConfig interface {
	getName() string
	getLinterNameForDiagnostic(*Diagnostic) string
	getAnalyzers() []*analysis.Analyzer
	useOriginalPackages() bool
	isTypecheckMode() bool
	reportIssues(*linter.Context) []Issue
	getLoadMode() LoadMode
}

func getIssuesCacheKey(analyzers []*analysis.Analyzer) string {
	return "lint/result:" + analyzersHashID(analyzers)
}

func saveIssuesToCache(allPkgs []*packages.Package, pkgsFromCache map[*packages.Package]bool,
	issues []result.Issue, lintCtx *linter.Context, analyzers []*analysis.Analyzer) {
	startedAt := time.Now()
	perPkgIssues := map[*packages.Package][]result.Issue{}
	for ind := range issues {
		i := &issues[ind]
		perPkgIssues[i.Pkg] = append(perPkgIssues[i.Pkg], *i)
	}

	savedIssuesCount := int32(0)
	lintResKey := getIssuesCacheKey(analyzers)

	workerCount := runtime.GOMAXPROCS(-1)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	pkgCh := make(chan *packages.Package, len(allPkgs))
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for pkg := range pkgCh {
				pkgIssues := perPkgIssues[pkg]
				encodedIssues := make([]EncodingIssue, 0, len(pkgIssues))
				for ind := range pkgIssues {
					i := &pkgIssues[ind]
					encodedIssues = append(encodedIssues, EncodingIssue{
						FromLinter:  i.FromLinter,
						Text:        i.Text,
						Pos:         i.Pos,
						LineRange:   i.LineRange,
						Replacement: i.Replacement,
					})
				}

				atomic.AddInt32(&savedIssuesCount, int32(len(encodedIssues)))
				if err := lintCtx.PkgCache.Put(pkg, pkgcache.HashModeNeedAllDeps, lintResKey, encodedIssues); err != nil {
					lintCtx.Log.Infof("Failed to save package %s issues (%d) to cache: %s", pkg, len(pkgIssues), err)
				} else {
					issuesCacheDebugf("Saved package %s issues (%d) to cache", pkg, len(pkgIssues))
				}
			}
		}()
	}

	for _, pkg := range allPkgs {
		if pkgsFromCache[pkg] {
			continue
		}

		pkgCh <- pkg
	}
	close(pkgCh)
	wg.Wait()

	issuesCacheDebugf("Saved %d issues from %d packages to cache in %s", savedIssuesCount, len(allPkgs), time.Since(startedAt))
}

//nolint:gocritic
func loadIssuesFromCache(pkgs []*packages.Package, lintCtx *linter.Context,
	analyzers []*analysis.Analyzer) ([]result.Issue, map[*packages.Package]bool) {
	startedAt := time.Now()

	lintResKey := getIssuesCacheKey(analyzers)
	type cacheRes struct {
		issues  []result.Issue
		loadErr error
	}
	pkgToCacheRes := make(map[*packages.Package]*cacheRes, len(pkgs))
	for _, pkg := range pkgs {
		pkgToCacheRes[pkg] = &cacheRes{}
	}

	workerCount := runtime.GOMAXPROCS(-1)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	pkgCh := make(chan *packages.Package, len(pkgs))
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for pkg := range pkgCh {
				var pkgIssues []EncodingIssue
				err := lintCtx.PkgCache.Get(pkg, pkgcache.HashModeNeedAllDeps, lintResKey, &pkgIssues)
				cacheRes := pkgToCacheRes[pkg]
				cacheRes.loadErr = err
				if err != nil {
					continue
				}
				if len(pkgIssues) == 0 {
					continue
				}

				issues := make([]result.Issue, 0, len(pkgIssues))
				for _, i := range pkgIssues {
					issues = append(issues, result.Issue{
						FromLinter:  i.FromLinter,
						Text:        i.Text,
						Pos:         i.Pos,
						LineRange:   i.LineRange,
						Replacement: i.Replacement,
						Pkg:         pkg,
					})
				}
				cacheRes.issues = issues
			}
		}()
	}

	for _, pkg := range pkgs {
		pkgCh <- pkg
	}
	close(pkgCh)
	wg.Wait()

	loadedIssuesCount := 0
	var issues []result.Issue
	pkgsFromCache := map[*packages.Package]bool{}
	for pkg, cacheRes := range pkgToCacheRes {
		if cacheRes.loadErr == nil {
			loadedIssuesCount += len(cacheRes.issues)
			pkgsFromCache[pkg] = true
			issues = append(issues, cacheRes.issues...)
			issuesCacheDebugf("Loaded package %s issues (%d) from cache", pkg, len(cacheRes.issues))
		} else {
			issuesCacheDebugf("Didn't load package %s issues from cache: %s", pkg, cacheRes.loadErr)
		}
	}
	issuesCacheDebugf("Loaded %d issues from cache in %s, analyzing %d/%d packages",
		loadedIssuesCount, time.Since(startedAt), len(pkgs)-len(pkgsFromCache), len(pkgs))
	return issues, pkgsFromCache
}

func runAnalyzers(cfg runAnalyzersConfig, lintCtx *linter.Context) ([]result.Issue, error) {
	log := lintCtx.Log.Child("goanalysis")
	sw := timeutils.NewStopwatch("analyzers", log)
	defer sw.PrintTopStages(10)

	runner := newRunner(cfg.getName(), log, lintCtx.PkgCache, lintCtx.LoadGuard, cfg.getLoadMode(), sw)

	pkgs := lintCtx.Packages
	if cfg.useOriginalPackages() {
		pkgs = lintCtx.OriginalPackages
	}

	issues, pkgsFromCache := loadIssuesFromCache(pkgs, lintCtx, cfg.getAnalyzers())
	var pkgsToAnalyze []*packages.Package
	for _, pkg := range pkgs {
		if !pkgsFromCache[pkg] {
			pkgsToAnalyze = append(pkgsToAnalyze, pkg)
		}
	}

	diags, errs, passToPkg := runner.run(cfg.getAnalyzers(), pkgsToAnalyze)

	defer func() {
		if len(errs) == 0 {
			// If we try to save to cache even if we have compilation errors
			// we won't see them on repeated runs.
			saveIssuesToCache(pkgs, pkgsFromCache, issues, lintCtx, cfg.getAnalyzers())
		}
	}()

	buildAllIssues := func() []result.Issue {
		var retIssues []result.Issue
		reportedIssues := cfg.reportIssues(lintCtx)
		for i := range reportedIssues {
			issue := &reportedIssues[i].Issue
			if issue.Pkg == nil {
				issue.Pkg = passToPkg[reportedIssues[i].Pass]
			}
			retIssues = append(retIssues, *issue)
		}
		retIssues = append(retIssues, buildIssues(diags, cfg.getLinterNameForDiagnostic)...)
		return retIssues
	}

	if cfg.isTypecheckMode() {
		errIssues, err := buildIssuesFromErrorsForTypecheckMode(errs, lintCtx)
		if err != nil {
			return nil, err
		}

		issues = append(issues, errIssues...)
		issues = append(issues, buildAllIssues()...)

		return issues, nil
	}

	// Don't print all errs: they can duplicate.
	if len(errs) != 0 {
		return nil, errs[0]
	}

	issues = append(issues, buildAllIssues()...)
	return issues, nil
}

func (lnt *Linter) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	if err := lnt.preRun(lintCtx); err != nil {
		return nil, err
	}

	return runAnalyzers(lnt, lintCtx)
}

func analyzersHashID(analyzers []*analysis.Analyzer) string {
	names := make([]string, 0, len(analyzers))
	for _, a := range analyzers {
		names = append(names, a.Name)
	}

	sort.Strings(names)
	return strings.Join(names, ",")
}
