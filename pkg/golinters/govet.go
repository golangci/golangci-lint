package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/goutils"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/timeutils"
	govetAPI "github.com/golangci/govet"
)

type Govet struct{}

func (Govet) Name() string {
	return "govet"
}

func (Govet) Desc() string {
	return "Vet examines Go source code and reports suspicious constructs, " +
		"such as Printf calls whose arguments do not align with the format string"
}

func (g Govet) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var govetIssues []govetAPI.Issue
	var err error
	if lintCtx.Settings().Govet.UseInstalledPackages {
		govetIssues, err = g.runOnInstalledPackages(ctx, lintCtx)
		if err != nil {
			return nil, fmt.Errorf("can't run govet on installed packages: %s", err)
		}
	} else {
		govetIssues, err = g.runOnSourcePackages(ctx, lintCtx)
		if err != nil {
			return nil, fmt.Errorf("can't run govet on source packages: %s", err)
		}
	}

	if len(govetIssues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(govetIssues))
	for _, i := range govetIssues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       i.Message,
			FromLinter: g.Name(),
		})
	}
	return res, nil
}

func (g Govet) runOnInstalledPackages(ctx context.Context, lintCtx *linter.Context) ([]govetAPI.Issue, error) {
	if err := g.installPackages(ctx, lintCtx); err != nil {
		return nil, fmt.Errorf("can't install packages (it's required for govet): %s", err)
	}

	// TODO: check .S asm files: govet can do it if pass dirs
	var govetIssues []govetAPI.Issue
	for _, pkg := range lintCtx.PkgProgram.Packages() {
		var astFiles []*ast.File
		var fset *token.FileSet
		for _, fname := range pkg.Files(lintCtx.Cfg.Run.AnalyzeTests) {
			af := lintCtx.ASTCache.Get(fname)
			if af == nil || af.Err != nil {
				return nil, fmt.Errorf("can't get parsed file %q from ast cache: %#v", fname, af)
			}
			astFiles = append(astFiles, af.F)
			fset = af.Fset
		}
		if len(astFiles) == 0 {
			continue
		}
		issues, err := govetAPI.Analyze(astFiles, fset, nil,
			lintCtx.Settings().Govet.CheckShadowing)
		if err != nil {
			return nil, err
		}
		govetIssues = append(govetIssues, issues...)
	}

	return govetIssues, nil
}

func (g Govet) installPackages(ctx context.Context, lintCtx *linter.Context) error {
	inGoRoot, err := goutils.InGoRoot()
	if err != nil {
		return fmt.Errorf("can't check whether we are in $GOROOT: %s", err)
	}

	if inGoRoot {
		// Go source packages already should be installed into $GOROOT/pkg with go distribution
		lintCtx.Log.Infof("In $GOROOT, don't install packages")
		return nil
	}

	if err := g.installNonTestPackages(ctx, lintCtx); err != nil {
		return err
	}

	if err := g.installTestDependencies(ctx, lintCtx); err != nil {
		return err
	}

	return nil
}

func (g Govet) installTestDependencies(ctx context.Context, lintCtx *linter.Context) error {
	log := lintCtx.Log
	packages := lintCtx.PkgProgram.Packages()
	var testDirs []string
	for _, pkg := range packages {
		dir := pkg.Dir()
		if dir == "" {
			log.Warnf("Package %#v has empty dir", pkg)
			continue
		}

		if !strings.HasPrefix(dir, ".") {
			// go install can't work without that
			dir = "./" + dir
		}

		if len(pkg.TestFiles()) != 0 {
			testDirs = append(testDirs, dir)
		}
	}

	if len(testDirs) == 0 {
		log.Infof("No test files in packages %#v", packages)
		return nil
	}

	args := append([]string{"test", "-i"}, testDirs...)
	return runGoCommand(ctx, log, args...)
}

func (g Govet) installNonTestPackages(ctx context.Context, lintCtx *linter.Context) error {
	log := lintCtx.Log
	packages := lintCtx.PkgProgram.Packages()
	var importPaths []string
	for _, pkg := range packages {
		if pkg.IsTestOnly() {
			// test-only package will be processed by installTestDependencies
			continue
		}

		dir := pkg.Dir()
		if dir == "" {
			log.Warnf("Package %#v has empty dir", pkg)
			continue
		}

		if !strings.HasPrefix(dir, ".") {
			// go install can't work without that
			dir = "./" + dir
		}

		importPaths = append(importPaths, dir)
	}

	if len(importPaths) == 0 {
		log.Infof("No packages to install, all packages: %#v", packages)
		return nil
	}

	// we need type information of dependencies of analyzed packages
	// so we pass -i option to install it
	if err := runGoInstall(ctx, log, importPaths, true); err != nil {
		// try without -i option: go < 1.10 doesn't support this option
		// and install dependencies by default.
		return runGoInstall(ctx, log, importPaths, false)
	}

	return nil
}

func runGoInstall(ctx context.Context, log logutils.Log, importPaths []string, withIOption bool) error {
	args := []string{"install"}
	if withIOption {
		args = append(args, "-i")
	}
	args = append(args, importPaths...)

	return runGoCommand(ctx, log, args...)
}

func runGoCommand(ctx context.Context, log logutils.Log, args ...string) error {
	argsStr := strings.Join(args, " ")
	defer timeutils.Track(time.Now(), log, "go %s", argsStr)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = append([]string{}, os.Environ()...)
	cmd.Env = append(cmd.Env, "GOMAXPROCS=1") // don't consume more than 1 cpu

	// use .Output but not .Run to capture StdErr in err
	_, err := cmd.Output()
	if err != nil {
		var stderr string
		if ee, ok := err.(*exec.ExitError); ok && ee.Stderr != nil {
			stderr = ": " + string(ee.Stderr)
		}

		return fmt.Errorf("can't run [go %s]: %s%s", argsStr, err, stderr)
	}

	return nil
}

func (g Govet) runOnSourcePackages(ctx context.Context, lintCtx *linter.Context) ([]govetAPI.Issue, error) {
	// TODO: check .S asm files: govet can do it if pass dirs
	var govetIssues []govetAPI.Issue
	for _, pkg := range lintCtx.Program.InitialPackages() {
		if len(pkg.Files) == 0 {
			continue
		}
		issues, err := govetAPI.Analyze(pkg.Files, lintCtx.Program.Fset, pkg,
			lintCtx.Settings().Govet.CheckShadowing)
		if err != nil {
			return nil, err
		}
		govetIssues = append(govetIssues, issues...)
	}

	return govetIssues, nil
}
