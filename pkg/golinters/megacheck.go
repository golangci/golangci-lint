package golinters

import (
	"context"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"honnef.co/go/tools/unused"

	"honnef.co/go/tools/lint"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	MegacheckParentName      = "megacheck"
	MegacheckStaticcheckName = "staticcheck"
	MegacheckUnusedName      = "unused"
	MegacheckGosimpleName    = "gosimple"
	MegacheckStylecheckName  = "stylecheck"
)

var debugf = logutils.Debug("megacheck")

type Staticcheck struct {
	megacheck
}

func NewStaticcheck() *Staticcheck {
	return &Staticcheck{
		megacheck: megacheck{
			staticcheckEnabled: true,
		},
	}
}

func (Staticcheck) Name() string { return MegacheckStaticcheckName }
func (Staticcheck) Desc() string {
	return "Staticcheck is a go vet on steroids, applying a ton of static analysis checks"
}

type Gosimple struct {
	megacheck
}

func NewGosimple() *Gosimple {
	return &Gosimple{
		megacheck: megacheck{
			gosimpleEnabled: true,
		},
	}
}

func (Gosimple) Name() string { return MegacheckGosimpleName }
func (Gosimple) Desc() string {
	return "Linter for Go source code that specializes in simplifying a code"
}

type Unused struct {
	megacheck
}

func NewUnused() *Unused {
	return &Unused{
		megacheck: megacheck{
			unusedEnabled: true,
		},
	}
}

func (Unused) Name() string { return MegacheckUnusedName }
func (Unused) Desc() string {
	return "Checks Go code for unused constants, variables, functions and types"
}

type Stylecheck struct {
	megacheck
}

func NewStylecheck() *Stylecheck {
	return &Stylecheck{
		megacheck: megacheck{
			stylecheckEnabled: true,
		},
	}
}

func (Stylecheck) Name() string { return MegacheckStylecheckName }
func (Stylecheck) Desc() string { return "Stylecheck is a replacement for golint" }

type megacheck struct {
	unusedEnabled      bool
	gosimpleEnabled    bool
	staticcheckEnabled bool
	stylecheckEnabled  bool
}

func (megacheck) Name() string {
	return MegacheckParentName
}

func (megacheck) Desc() string {
	return "" // shouldn't be called
}

func (m *megacheck) enableChildLinter(name string) error {
	switch name {
	case MegacheckStaticcheckName:
		m.staticcheckEnabled = true
	case MegacheckGosimpleName:
		m.gosimpleEnabled = true
	case MegacheckUnusedName:
		m.unusedEnabled = true
	case MegacheckStylecheckName:
		m.stylecheckEnabled = true
	default:
		return fmt.Errorf("invalid child linter name %s for metalinter %s", name, m.Name())
	}

	return nil
}

type MegacheckMetalinter struct{}

func (MegacheckMetalinter) Name() string {
	return MegacheckParentName
}

func (MegacheckMetalinter) BuildLinterConfig(enabledChildren []string) (*linter.Config, error) {
	var m megacheck
	for _, name := range enabledChildren {
		if err := m.enableChildLinter(name); err != nil {
			return nil, err
		}
	}

	// TODO: merge linter.Config and linter.Linter or refactor it in another way
	lc := &linter.Config{
		Linter:           m,
		EnabledByDefault: false,
		NeedsSSARepr:     false,
		InPresets:        []string{linter.PresetStyle, linter.PresetBugs, linter.PresetUnused},
		Speed:            1,
		AlternativeNames: nil,
		OriginalURL:      "",
		ParentLinterName: "",
	}
	if m.unusedEnabled {
		lc = lc.WithLoadDepsTypeInfo()
	} else {
		lc = lc.WithLoadForGoAnalysis()
	}
	return lc, nil
}

func (MegacheckMetalinter) DefaultChildLinterNames() []string {
	// no stylecheck here for backwards compatibility for users who enabled megacheck: don't enable extra
	// linter for them
	return []string{MegacheckStaticcheckName, MegacheckGosimpleName, MegacheckUnusedName}
}

func (m MegacheckMetalinter) AllChildLinterNames() []string {
	return append(m.DefaultChildLinterNames(), MegacheckStylecheckName)
}

func (m megacheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	// Use OriginalPackages not Packages because `unused` doesn't work properly
	// when we deduplicate normal and test packages.
	return m.runMegacheck(ctx, lintCtx)
}

func getAnalyzers(m map[string]*analysis.Analyzer) []*analysis.Analyzer {
	var ret []*analysis.Analyzer
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func setGoVersion(analyzers []*analysis.Analyzer) {
	const goVersion = 13 // TODO
	for _, a := range analyzers {
		if v := a.Flags.Lookup("go"); v != nil {
			if err := v.Value.Set(fmt.Sprintf("1.%d", goVersion)); err != nil {
				debugf("Failed to set go version: %s", err)
			}
		}
	}
}

func (m megacheck) runMegacheck(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var linters []linter.Linter

	if m.gosimpleEnabled {
		analyzers := getAnalyzers(simple.Analyzers)
		setGoVersion(analyzers)
		lnt := goanalysis.NewLinter(MegacheckGosimpleName, "", analyzers, nil)
		linters = append(linters, lnt)
	}
	if m.staticcheckEnabled {
		analyzers := getAnalyzers(staticcheck.Analyzers)
		setGoVersion(analyzers)
		lnt := goanalysis.NewLinter(MegacheckStaticcheckName, "", analyzers, nil)
		linters = append(linters, lnt)
	}
	if m.stylecheckEnabled {
		analyzers := getAnalyzers(stylecheck.Analyzers)
		setGoVersion(analyzers)
		lnt := goanalysis.NewLinter(MegacheckStylecheckName, "", analyzers, nil)
		linters = append(linters, lnt)
	}

	var u lint.CumulativeChecker
	if m.unusedEnabled {
		u = unused.NewChecker(lintCtx.Settings().Unused.CheckExported)
		analyzers := []*analysis.Analyzer{u.Analyzer()}
		setGoVersion(analyzers)
		lnt := goanalysis.NewLinter(MegacheckUnusedName, "", analyzers, nil)
		linters = append(linters, lnt)
	}

	if len(linters) == 0 {
		return nil, nil
	}

	var issues []result.Issue
	for _, lnt := range linters {
		i, err := lnt.Run(ctx, lintCtx)
		if err != nil {
			return nil, err
		}
		issues = append(issues, i...)
	}

	if u != nil {
		for _, ur := range u.Result() {
			p := u.ProblemObject(lintCtx.Packages[0].Fset, ur)
			issues = append(issues, result.Issue{
				FromLinter: MegacheckUnusedName,
				Text:       p.Message,
				Pos:        p.Pos,
			})
		}
	}

	return issues, nil
}

func (m megacheck) Analyzers() []*analysis.Analyzer {
	if m.unusedEnabled {
		// Don't treat this linter as go/analysis linter if unused is used
		// because it has non-standard API.
		return nil
	}

	var allAnalyzers []*analysis.Analyzer
	if m.gosimpleEnabled {
		allAnalyzers = append(allAnalyzers, getAnalyzers(simple.Analyzers)...)
	}
	if m.staticcheckEnabled {
		allAnalyzers = append(allAnalyzers, getAnalyzers(staticcheck.Analyzers)...)
	}
	if m.stylecheckEnabled {
		allAnalyzers = append(allAnalyzers, getAnalyzers(stylecheck.Analyzers)...)
	}
	setGoVersion(allAnalyzers)
	return allAnalyzers
}

func (megacheck) Cfg() map[string]map[string]interface{} {
	return nil
}

func (m megacheck) AnalyzerToLinterNameMapping() map[*analysis.Analyzer]string {
	ret := map[*analysis.Analyzer]string{}
	if m.gosimpleEnabled {
		for _, a := range simple.Analyzers {
			ret[a] = MegacheckGosimpleName
		}
	}
	if m.staticcheckEnabled {
		for _, a := range staticcheck.Analyzers {
			ret[a] = MegacheckStaticcheckName
		}
	}
	if m.stylecheckEnabled {
		for _, a := range stylecheck.Analyzers {
			ret[a] = MegacheckStylecheckName
		}
	}
	return ret
}
