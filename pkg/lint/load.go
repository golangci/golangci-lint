package lint

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/ssa/ssautil"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
)

type Context struct {
	Paths                *fsutils.ProjectPaths
	Cfg                  *config.Config
	Program              *loader.Program
	SSAProgram           *ssa.Program
	LoaderConfig         *loader.Config
	ASTCache             *astcache.Cache
	NotCompilingPackages []*loader.PackageInfo
}

func (c *Context) Settings() *config.LintersSettings {
	return &c.Cfg.LintersSettings
}

func isFullImportNeeded(linters []linter.Config) bool {
	for _, linter := range linters {
		if linter.NeedsProgramLoading() {
			return true
		}
	}

	return false
}

func isSSAReprNeeded(linters []linter.Config) bool {
	for _, linter := range linters {
		if linter.NeedsSSARepresentation() {
			return true
		}
	}

	return false
}

func loadWholeAppIfNeeded(ctx context.Context, linters []linter.Config, cfg *config.Run, paths *fsutils.ProjectPaths) (*loader.Program, *loader.Config, error) {
	if !isFullImportNeeded(linters) {
		return nil, nil, nil
	}

	startedAt := time.Now()
	defer func() {
		logrus.Infof("Program loading took %s", time.Since(startedAt))
	}()

	bctx := build.Default
	bctx.BuildTags = append(bctx.BuildTags, cfg.BuildTags...)
	loadcfg := &loader.Config{
		Build:       &bctx,
		AllowErrors: true, // Try to analyze event partially
	}
	rest, err := loadcfg.FromArgs(paths.MixedPaths(), cfg.AnalyzeTests)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parepare load config with paths: %s", err)
	}
	if len(rest) > 0 {
		return nil, nil, fmt.Errorf("unhandled loading paths: %v", rest)
	}

	prog, err := loadcfg.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("can't load program from paths %v: %s", paths.MixedPaths(), err)
	}

	return prog, loadcfg, nil
}

func buildSSAProgram(ctx context.Context, lprog *loader.Program) *ssa.Program {
	startedAt := time.Now()
	defer func() {
		logrus.Infof("SSA repr building took %s", time.Since(startedAt))
	}()

	ssaProg := ssautil.CreateProgram(lprog, ssa.GlobalDebug)
	ssaProg.Build()
	return ssaProg
}

func discoverGoRoot() (string, error) {
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		return goroot, nil
	}

	output, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return "", fmt.Errorf("can't execute go env GOROOT: %s", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// separateNotCompilingPackages moves not compiling packages into separate slices:
// a lot of linters crash on such packages. Leave them only for those linters
// which can work with them.
func separateNotCompilingPackages(lintCtx *linter.Context) {
	prog := lintCtx.Program

	if prog.Created != nil {
		compilingCreated := make([]*loader.PackageInfo, 0, len(prog.Created))
		for _, info := range prog.Created {
			if len(info.Errors) != 0 {
				lintCtx.NotCompilingPackages = append(lintCtx.NotCompilingPackages, info)
			} else {
				compilingCreated = append(compilingCreated, info)
			}
		}
		prog.Created = compilingCreated
	}

	if prog.Imported != nil {
		for k, info := range prog.Imported {
			if len(info.Errors) != 0 {
				lintCtx.NotCompilingPackages = append(lintCtx.NotCompilingPackages, info)
				delete(prog.Imported, k)
			}
		}
	}
}

func LoadContext(ctx context.Context, linters []linter.Config, cfg *config.Config) (*linter.Context, error) {
	// Set GOROOT to have working cross-compilation: cross-compiled binaries
	// have invalid GOROOT. XXX: can't use runtime.GOROOT().
	goroot, err := discoverGoRoot()
	if err != nil {
		return nil, fmt.Errorf("can't discover GOROOT: %s", err)
	}
	os.Setenv("GOROOT", goroot)
	build.Default.GOROOT = goroot
	logrus.Infof("set GOROOT=%q", goroot)

	args := cfg.Run.Args
	if len(args) == 0 {
		args = []string{"./..."}
	}

	paths, err := fsutils.GetPathsForAnalysis(ctx, args, cfg.Run.AnalyzeTests)
	if err != nil {
		return nil, err
	}

	prog, loaderConfig, err := loadWholeAppIfNeeded(ctx, linters, &cfg.Run, paths)
	if err != nil {
		return nil, err
	}

	var ssaProg *ssa.Program
	if prog != nil && isSSAReprNeeded(linters) {
		ssaProg = buildSSAProgram(ctx, prog)
	}

	var astCache *astcache.Cache
	if prog != nil {
		astCache, err = astcache.LoadFromProgram(prog)
		if err != nil {
			return nil, err
		}
	} else {
		astCache = astcache.LoadFromFiles(paths.Files)
	}

	ret := &linter.Context{
		Paths:        paths,
		Cfg:          cfg,
		Program:      prog,
		SSAProgram:   ssaProg,
		LoaderConfig: loaderConfig,
		ASTCache:     astCache,
	}

	if prog != nil {
		separateNotCompilingPackages(ret)
	}

	return ret, nil
}
