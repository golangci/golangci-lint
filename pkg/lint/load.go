package lint

import (
	"context"
	"fmt"
	"go/build"
	"go/parser"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/ssa/ssautil"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
)

var loadDebugf = logutils.Debug("load")

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

func normalizePaths(paths []string) ([]string, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working dir: %s", err)
	}

	ret := make([]string, 0, len(paths))
	for _, p := range paths {
		if filepath.IsAbs(p) {
			relPath, err := filepath.Rel(root, p)
			if err != nil {
				return nil, fmt.Errorf("can't get relative path for path %s and root %s: %s",
					p, root, err)
			}
			p = relPath
		}

		ret = append(ret, "./"+p)
	}

	return ret, nil
}

func getCurrentProjectImportPath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("no GOPATH env variable")
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("can't get workind directory: %s", err)
	}

	if !strings.HasPrefix(wd, gopath) {
		return "", fmt.Errorf("currently no in gopath: %q isn't a prefix of %q", gopath, wd)
	}

	path := strings.TrimPrefix(wd, gopath)
	path = strings.TrimPrefix(path, string(os.PathSeparator)) // if GOPATH contains separator at the end
	src := "src" + string(os.PathSeparator)
	if !strings.HasPrefix(path, src) {
		return "", fmt.Errorf("currently no in gopath/src: %q isn't a prefix of %q", src, path)
	}

	path = strings.TrimPrefix(path, src)
	path = strings.Replace(path, string(os.PathSeparator), "/", -1)
	return path, nil
}

func isLocalProjectAnalysis(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "..") || filepath.IsAbs(arg) {
			return false
		}
	}

	return true
}

func getTypeCheckFuncBodies(cfg *config.Run, linters []linter.Config, pkgProg *packages.Program) func(string) bool {
	if !isLocalProjectAnalysis(cfg.Args) {
		loadDebugf("analysis in nonlocal, don't optimize loading by not typechecking func bodies")
		return nil
	}

	if isSSAReprNeeded(linters) {
		loadDebugf("ssa repr is needed, don't optimize loading by not typechecking func bodies")
		return nil
	}

	if len(pkgProg.Dirs()) == 0 {
		// files run, in this mode packages are fake: can't check their path properly
		return nil
	}

	projPath, err := getCurrentProjectImportPath()
	if err != nil {
		logrus.Infof("can't get cur project path: %s", err)
		return nil
	}

	return func(path string) bool {
		if strings.HasPrefix(path, ".") {
			loadDebugf("%s: dot import: typecheck func bodies", path)
			return true
		}

		isLocalPath := strings.HasPrefix(path, projPath)
		if isLocalPath {
			localPath := strings.TrimPrefix(path, projPath)
			localPath = strings.TrimPrefix(localPath, "/")
			if strings.HasPrefix(localPath, "vendor/") {
				loadDebugf("%s: local vendor import: DO NOT typecheck func bodies", path)
				return false
			}

			loadDebugf("%s: local import: typecheck func bodies", path)
			return true
		}

		loadDebugf("%s: not local import: DO NOT typecheck func bodies", path)
		return false
	}
}

func loadWholeAppIfNeeded(ctx context.Context, linters []linter.Config, cfg *config.Run, pkgProg *packages.Program) (*loader.Program, *loader.Config, error) {
	if !isFullImportNeeded(linters) {
		return nil, nil, nil
	}

	startedAt := time.Now()
	defer func() {
		logrus.Infof("Program loading took %s", time.Since(startedAt))
	}()

	bctx := pkgProg.BuildContext()
	loadcfg := &loader.Config{
		Build:               bctx,
		AllowErrors:         true,                 // Try to analyze partially
		ParserMode:          parser.ParseComments, // AST will be reused by linters
		TypeCheckFuncBodies: getTypeCheckFuncBodies(cfg, linters, pkgProg),
	}

	var loaderArgs []string
	dirs := pkgProg.Dirs()
	if len(dirs) != 0 {
		loaderArgs = dirs // dirs run
	} else {
		loaderArgs = pkgProg.Files(cfg.AnalyzeTests) // files run
	}

	nLoaderArgs, err := normalizePaths(loaderArgs)
	if err != nil {
		return nil, nil, err
	}

	rest, err := loadcfg.FromArgs(nLoaderArgs, cfg.AnalyzeTests)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parepare load config with paths: %s", err)
	}
	if len(rest) > 0 {
		return nil, nil, fmt.Errorf("unhandled loading paths: %v", rest)
	}

	prog, err := loadcfg.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("can't load program from paths %v: %s", loaderArgs, err)
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

//nolint:gocyclo
func LoadContext(ctx context.Context, linters []linter.Config, cfg *config.Config) (*linter.Context, error) {
	// Set GOROOT to have working cross-compilation: cross-compiled binaries
	// have invalid GOROOT. XXX: can't use runtime.GOROOT().
	goroot, err := discoverGoRoot()
	if err != nil {
		return nil, fmt.Errorf("can't discover GOROOT: %s", err)
	}
	os.Setenv("GOROOT", goroot)
	build.Default.GOROOT = goroot

	args := cfg.Run.Args
	if len(args) == 0 {
		args = []string{"./..."}
	}

	skipDirs := append([]string{}, packages.StdExcludeDirRegexps...)
	skipDirs = append(skipDirs, cfg.Run.SkipDirs...)
	r, err := packages.NewResolver(cfg.Run.BuildTags, skipDirs)
	if err != nil {
		return nil, err
	}

	pkgProg, err := r.Resolve(args...)
	if err != nil {
		return nil, err
	}

	prog, loaderConfig, err := loadWholeAppIfNeeded(ctx, linters, &cfg.Run, pkgProg)
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
	} else {
		astCache, err = astcache.LoadFromFiles(pkgProg.Files(cfg.Run.AnalyzeTests))
	}
	if err != nil {
		return nil, err
	}

	ret := &linter.Context{
		PkgProgram:   pkgProg,
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
