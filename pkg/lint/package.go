package lint

import (
	"context"
	"fmt"
	"go/build"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ldez/grignotin/goenv"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/exitcodes"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis/load"
	"github.com/golangci/golangci-lint/v2/pkg/goutil"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

// PackageLoader loads packages based on [golang.org/x/tools/go/packages.Load].
type PackageLoader struct {
	log    logutils.Log
	debugf logutils.DebugFunc

	cfg *config.Config

	args []string

	pkgTestIDRe *regexp.Regexp

	goenv *goutil.Env

	loadGuard *load.Guard
}

// NewPackageLoader creates a new PackageLoader.
func NewPackageLoader(log logutils.Log, cfg *config.Config, args []string, env *goutil.Env, loadGuard *load.Guard) *PackageLoader {
	return &PackageLoader{
		cfg:         cfg,
		args:        args,
		log:         log,
		debugf:      logutils.Debug(logutils.DebugKeyLoader),
		goenv:       env,
		pkgTestIDRe: regexp.MustCompile(`^(.*) \[(.*)\.test\]`),
		loadGuard:   loadGuard,
	}
}

// Load loads packages.
func (l *PackageLoader) Load(ctx context.Context, linters []*linter.Config) (pkgs, deduplicatedPkgs []*packages.Package, err error) {
	// Check for multiple modules and provide helpful error
	if err := l.detectMultipleModules(ctx); err != nil {
		return nil, nil, err
	}

	loadMode := findLoadMode(linters)

	pkgs, loadErr := l.loadPackages(ctx, loadMode)
	if loadErr != nil {
		return nil, nil, fmt.Errorf("failed to load packages: %w", loadErr)
	}

	return pkgs, l.filterDuplicatePackages(pkgs), nil
}

// detectMultipleModules checks if multiple arguments refer to different modules
func (l *PackageLoader) detectMultipleModules(ctx context.Context) error {
	if len(l.args) <= 1 {
		return nil
	}

	var moduleRoots []string
	seenRoots := make(map[string]bool)

	for _, arg := range l.args {
		moduleRoot := l.findModuleRootForArg(ctx, arg)
		if moduleRoot != "" && !seenRoots[moduleRoot] {
			moduleRoots = append(moduleRoots, moduleRoot)
			seenRoots[moduleRoot] = true
		}
	}

	if len(moduleRoots) > 1 {
		return fmt.Errorf("multiple Go modules detected: %v\n\n"+
			"Multi-module analysis is not supported. Each module should be analyzed separately:\n"+
			"  golangci-lint run %s\n  golangci-lint run %s",
			moduleRoots, moduleRoots[0], moduleRoots[1])
	}

	return nil
}

// findModuleRootForArg finds the module root for a given argument using go env
func (l *PackageLoader) findModuleRootForArg(ctx context.Context, arg string) string {
	absPath, err := filepath.Abs(arg)
	if err != nil {
		if l.debugf != nil {
			l.debugf("Failed to get absolute path for %s: %v", arg, err)
		}
		return ""
	}

	// Determine the directory to check
	var targetDir string
	if info, statErr := os.Stat(absPath); statErr == nil && info.IsDir() {
		targetDir = absPath
	} else if statErr == nil {
		targetDir = filepath.Dir(absPath)
	} else {
		return ""
	}

	// Save current directory
	originalWd, err := os.Getwd()
	if err != nil {
		if l.debugf != nil {
			l.debugf("Failed to get current directory: %v", err)
		}
		return ""
	}
	defer func() {
		if chErr := os.Chdir(originalWd); chErr != nil && l.debugf != nil {
			l.debugf("Failed to restore directory %s: %v", originalWd, chErr)
		}
	}()

	// Change to target directory and use go env GOMOD
	if chdirErr := os.Chdir(targetDir); chdirErr != nil {
		if l.debugf != nil {
			l.debugf("Failed to change to directory %s: %v", targetDir, chdirErr)
		}
		return ""
	}

	goModPath, err := goenv.GetOne(ctx, goenv.GOMOD)
	if err != nil || goModPath == "" {
		if l.debugf != nil {
			l.debugf("go env GOMOD failed in %s: err=%v, path=%s", targetDir, err, goModPath)
		}
		return ""
	}

	return filepath.Dir(goModPath)
}

// determineWorkingDir determines the working directory for package loading.
// If the first argument is within a directory tree that has a go.mod file, use that module root.
// Otherwise, use the current working directory.
func (l *PackageLoader) determineWorkingDir(ctx context.Context) string {
	if len(l.args) == 0 {
		return ""
	}

	moduleRoot := l.findModuleRootForArg(ctx, l.args[0])
	if moduleRoot != "" {
		if l.debugf != nil {
			l.debugf("Found module root %s, using as working dir", moduleRoot)
		}
	}
	return moduleRoot
}

func (l *PackageLoader) loadPackages(ctx context.Context, loadMode packages.LoadMode) ([]*packages.Package, error) {
	defer func(startedAt time.Time) {
		l.log.Infof("Go packages loading at mode %s took %s", stringifyLoadMode(loadMode), time.Since(startedAt))
	}(time.Now())

	l.prepareBuildContext()

	conf := &packages.Config{
		Mode:       loadMode,
		Tests:      l.cfg.Run.AnalyzeTests,
		Context:    ctx,
		BuildFlags: l.makeBuildFlags(),
		Logf:       l.debugf,
		Dir:        l.determineWorkingDir(ctx),
		// TODO: use fset, parsefile, overlay
	}

	args := l.buildArgs(ctx)

	l.debugf("Built loader args are %s", args)

	pkgs, err := packages.Load(conf, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to load with go/packages: %w", err)
	}

	if loadMode&packages.NeedSyntax == 0 {
		// Needed e.g. for go/analysis loading.
		fset := token.NewFileSet()
		packages.Visit(pkgs, nil, func(pkg *packages.Package) {
			pkg.Fset = fset
			l.loadGuard.AddMutexForPkg(pkg)
		})
	}

	l.debugPrintLoadedPackages(pkgs)

	if err := l.parseLoadedPackagesErrors(pkgs); err != nil {
		return nil, err
	}

	return l.filterTestMainPackages(pkgs), nil
}

func (*PackageLoader) parseLoadedPackagesErrors(pkgs []*packages.Package) error {
	for _, pkg := range pkgs {
		var errs []packages.Error
		for _, err := range pkg.Errors {
			// quick fix: skip error related to `go list` invocation by packages.Load()
			// The behavior has been changed between go1.19 and go1.20, the error is now inside the JSON content.
			// https://github.com/golangci/golangci-lint/pull/3414#issuecomment-1364756303
			if strings.Contains(err.Msg, "# command-line-arguments") {
				continue
			}

			errs = append(errs, err)

			if strings.Contains(err.Msg, "no Go files") {
				return fmt.Errorf("package %s: %w", pkg.PkgPath, exitcodes.ErrNoGoFiles)
			}
			if strings.Contains(err.Msg, "cannot find package") {
				// when analyzing not existing directory
				return fmt.Errorf("%v: %w", err.Msg, exitcodes.ErrFailure)
			}
		}

		pkg.Errors = errs
	}

	return nil
}

func (l *PackageLoader) tryParseTestPackage(pkg *packages.Package) (name string, isTest bool) {
	matches := l.pkgTestIDRe.FindStringSubmatch(pkg.ID)
	if matches == nil {
		return "", false
	}

	return matches[1], true
}

func (l *PackageLoader) filterDuplicatePackages(pkgs []*packages.Package) []*packages.Package {
	packagesWithTests := map[string]bool{}
	for _, pkg := range pkgs {
		name, isTest := l.tryParseTestPackage(pkg)
		if !isTest {
			continue
		}
		packagesWithTests[name] = true
	}

	l.debugf("package with tests: %#v", packagesWithTests)

	var retPkgs []*packages.Package
	for _, pkg := range pkgs {
		_, isTest := l.tryParseTestPackage(pkg)
		if !isTest && packagesWithTests[pkg.PkgPath] {
			// If tests loading is enabled,
			// for package with files a.go and a_test.go go/packages loads two packages:
			// 1. ID=".../a" GoFiles=[a.go]
			// 2. ID=".../a [.../a.test]" GoFiles=[a.go a_test.go]
			// We need only the second package, otherwise we can get warnings about unused variables/fields/functions
			// in a.go if they are used only in a_test.go.
			l.debugf("skip pkg ID=%s because we load it with test package", pkg.ID)
			continue
		}

		retPkgs = append(retPkgs, pkg)
	}

	return retPkgs
}

func (l *PackageLoader) filterTestMainPackages(pkgs []*packages.Package) []*packages.Package {
	var retPkgs []*packages.Package
	for _, pkg := range pkgs {
		if pkg.Name == "main" && strings.HasSuffix(pkg.PkgPath, ".test") {
			// it's an implicit testmain package
			l.debugf("skip pkg ID=%s", pkg.ID)
			continue
		}

		retPkgs = append(retPkgs, pkg)
	}

	return retPkgs
}

func (l *PackageLoader) debugPrintLoadedPackages(pkgs []*packages.Package) {
	l.debugf("loaded %d pkgs", len(pkgs))
	for i, pkg := range pkgs {
		var syntaxFiles []string
		for _, sf := range pkg.Syntax {
			syntaxFiles = append(syntaxFiles, pkg.Fset.Position(sf.Pos()).Filename)
		}
		l.debugf("Loaded pkg #%d: ID=%s GoFiles=%s CompiledGoFiles=%s Syntax=%s",
			i, pkg.ID, pkg.GoFiles, pkg.CompiledGoFiles, syntaxFiles)
	}
}

func (l *PackageLoader) prepareBuildContext() {
	// Set GOROOT to have working cross-compilation: cross-compiled binaries
	// have invalid GOROOT. XXX: can't use runtime.GOROOT().
	goroot := l.goenv.Get(goenv.GOROOT)
	if goroot == "" {
		return
	}

	_ = os.Setenv(goenv.GOROOT, goroot)

	build.Default.GOROOT = goroot
	build.Default.BuildTags = l.cfg.Run.BuildTags
}

func (l *PackageLoader) makeBuildFlags() []string {
	var buildFlags []string

	if len(l.cfg.Run.BuildTags) != 0 {
		// go help build
		buildFlags = append(buildFlags, "-tags", strings.Join(l.cfg.Run.BuildTags, " "))
		l.log.Infof("Using build tags: %v", l.cfg.Run.BuildTags)
	}

	if l.cfg.Run.ModulesDownloadMode != "" {
		// go help modules
		buildFlags = append(buildFlags, fmt.Sprintf("-mod=%s", l.cfg.Run.ModulesDownloadMode))
	}

	return buildFlags
}

// buildArgs processes the arguments for package loading, handling directory changes appropriately.
func (l *PackageLoader) buildArgs(ctx context.Context) []string {
	if len(l.args) == 0 {
		return []string{"./..."}
	}

	workingDir := l.determineWorkingDir(ctx)

	// If we're using a different working directory, we need to adjust the arguments
	if workingDir != "" {
		// We're switching to the target directory as working dir, so use "./..." to analyze it
		return []string{"./..."}
	}

	// Use the original buildArgs logic for the current working directory
	var retArgs []string
	for _, arg := range l.args {
		if strings.HasPrefix(arg, ".") || filepath.IsAbs(arg) {
			retArgs = append(retArgs, arg)
		} else {
			// go/packages doesn't work well if we don't have the prefix ./ for local packages
			retArgs = append(retArgs, fmt.Sprintf(".%c%s", filepath.Separator, arg))
		}
	}

	return retArgs
}

func findLoadMode(linters []*linter.Config) packages.LoadMode {
	loadMode := packages.LoadMode(0)
	for _, lc := range linters {
		loadMode |= lc.LoadMode
	}

	return loadMode
}

func stringifyLoadMode(mode packages.LoadMode) string {
	m := map[packages.LoadMode]string{
		packages.NeedCompiledGoFiles: "compiled_files",
		packages.NeedDeps:            "deps",
		packages.NeedExportFile:      "exports_file",
		packages.NeedFiles:           "files",
		packages.NeedImports:         "imports",
		packages.NeedName:            "name",
		packages.NeedSyntax:          "syntax",
		packages.NeedTypes:           "types",
		packages.NeedTypesInfo:       "types_info",
		packages.NeedTypesSizes:      "types_sizes",
	}

	var flags []string
	for flag, flagStr := range m {
		if mode&flag != 0 {
			flags = append(flags, flagStr)
		}
	}

	return fmt.Sprintf("%d (%s)", mode, strings.Join(flags, "|"))
}
