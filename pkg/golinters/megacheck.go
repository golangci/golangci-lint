package golinters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"honnef.co/go/tools/lint"
	"honnef.co/go/tools/lint/lintutil"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/unused"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Megacheck struct {
	UnusedEnabled      bool
	GosimpleEnabled    bool
	StaticcheckEnabled bool
}

func (m Megacheck) Name() string {
	names := []string{}
	if m.UnusedEnabled {
		names = append(names, "unused")
	}
	if m.GosimpleEnabled {
		names = append(names, "gosimple")
	}
	if m.StaticcheckEnabled {
		names = append(names, "staticcheck")
	}

	if len(names) == 1 {
		return names[0] // only one sublinter is enabled
	}

	if len(names) == 3 {
		return "megacheck" // all enabled
	}

	return fmt.Sprintf("megacheck.{%s}", strings.Join(names, ","))
}

func (m Megacheck) Desc() string {
	descs := map[string]string{
		"unused":      "Checks Go code for unused constants, variables, functions and types",
		"gosimple":    "Linter for Go source code that specializes in simplifying a code",
		"staticcheck": "Staticcheck is a go vet on steroids, applying a ton of static analysis checks",
		"megacheck":   "3 sub-linters in one: unused, gosimple and staticcheck",
	}

	return descs[m.Name()]
}

func prettifyCompilationError(err error) error {
	i, _ := TypeCheck{}.parseError(err)
	if i == nil {
		return err
	}

	shortFilename, pathErr := fsutils.ShortestRelPath(i.Pos.Filename, "")
	if pathErr != nil {
		return err
	}

	errText := shortFilename
	if i.Line() != 0 {
		errText += fmt.Sprintf(":%d", i.Line())
	}
	errText += fmt.Sprintf(": %s", i.Text)
	return errors.New(errText)
}

func (m Megacheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	if len(lintCtx.NotCompilingPackages) != 0 {
		var packages []string
		var errors []error
		for _, p := range lintCtx.NotCompilingPackages {
			packages = append(packages, p.String())
			errors = append(errors, p.Errors...)
		}

		warnText := fmt.Sprintf("Can't run megacheck because of compilation errors in packages %s",
			packages)
		if len(errors) != 0 {
			warnText += fmt.Sprintf(": %s", prettifyCompilationError(errors[0]))
			if len(errors) > 1 {
				const runCmd = "golangci-lint run --no-config --disable-all -E typecheck"
				warnText += fmt.Sprintf(" and %d more errors: run `%s` to see all errors", len(errors)-1, runCmd)
			}
		}
		lintCtx.Log.Warnf("%s", warnText)

		// megacheck crashes if there are not compiling packages
		return nil, nil
	}

	issues := runMegacheck(lintCtx.Program, lintCtx.SSAProgram, lintCtx.LoaderConfig,
		m.StaticcheckEnabled, m.GosimpleEnabled, m.UnusedEnabled, lintCtx.Settings().Unused.CheckExported)
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Position,
			Text:       i.Text,
			FromLinter: m.Name(),
		})
	}
	return res, nil
}

func runMegacheck(program *loader.Program, ssaProg *ssa.Program, conf *loader.Config,
	enableStaticcheck, enableGosimple, enableUnused, checkExportedUnused bool) []lint.Problem {

	var checkers []lintutil.CheckerConfig

	if enableStaticcheck {
		sac := staticcheck.NewChecker()
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker: sac,
		})
	}

	if enableGosimple {
		sc := simple.NewChecker()
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker: sc,
		})
	}

	if enableUnused {
		uc := unused.NewChecker(unused.CheckAll)
		uc.WholeProgram = checkExportedUnused
		uc.ConsiderReflection = true
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker: unused.NewLintChecker(uc),
		})
	}

	fs := lintutil.FlagSet("megacheck")
	return lintutil.ProcessFlagSet(checkers, fs, program, ssaProg, conf)
}
