package golinters

import (
	"context"
	"fmt"

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
	return &linter.Config{
		Linter:            m,
		EnabledByDefault:  false,
		NeedsTypeInfo:     true,
		NeedsDepsTypeInfo: true,
		NeedsSSARepr:      false,
		InPresets:         []string{linter.PresetStyle, linter.PresetBugs, linter.PresetUnused},
		Speed:             1,
		AlternativeNames:  nil,
		OriginalURL:       "",
		ParentLinterName:  "",
	}, nil
}

func (MegacheckMetalinter) DefaultChildLinterNames() []string {
	// no stylecheck here for backwards compatibility for users who enabled megacheck: don't enable extra
	// linter for them
	return []string{MegacheckStaticcheckName, MegacheckGosimpleName, MegacheckUnusedName}
}

func (m MegacheckMetalinter) AllChildLinterNames() []string {
	return append(m.DefaultChildLinterNames(), MegacheckStylecheckName)
}

func (m MegacheckMetalinter) isValidChild(name string) bool {
	for _, child := range m.AllChildLinterNames() {
		if child == name {
			return true
		}
	}

	return false
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

func (m megacheck) runMegacheck(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var linters []linter.Linter

	if m.gosimpleEnabled {
		lnt := goanalysis.NewLinter(MegacheckGosimpleName, "", getAnalyzers(simple.Analyzers), nil)
		linters = append(linters, lnt)
	}
	if m.staticcheckEnabled {
		lnt := goanalysis.NewLinter(MegacheckStaticcheckName, "", getAnalyzers(staticcheck.Analyzers), nil)
		linters = append(linters, lnt)
	}
	if m.stylecheckEnabled {
		lnt := goanalysis.NewLinter(MegacheckStylecheckName, "", getAnalyzers(stylecheck.Analyzers), nil)
		linters = append(linters, lnt)
	}

	var u lint.CumulativeChecker
	if m.unusedEnabled {
		u = unused.NewChecker(lintCtx.Settings().Unused.CheckExported)
		lnt := goanalysis.NewLinter(MegacheckStylecheckName, "", []*analysis.Analyzer{u.Analyzer()}, nil)
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

	for _, ur := range u.Result() {
		p := u.ProblemObject(lintCtx.Packages[0].Fset, ur)
		issues = append(issues, result.Issue{
			FromLinter: MegacheckUnusedName,
			Text:       p.Message,
			Pos:        p.Pos,
		})
	}

	return issues, nil
}
