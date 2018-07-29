package packages

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Resolver struct {
	excludeDirs map[string]*regexp.Regexp
	buildTags   []string

	skippedDirs []string
	log         logutils.Log

	wd                  string // working directory
	importErrorsOccured int    // count of errors because too bad files in packages
}

func NewResolver(buildTags, excludeDirs []string, log logutils.Log) (*Resolver, error) {
	excludeDirsMap := map[string]*regexp.Regexp{}
	for _, dir := range excludeDirs {
		re, err := regexp.Compile(dir)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp %q: %s", dir, err)
		}

		excludeDirsMap[dir] = re
	}

	wd, err := fsutils.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working dir: %s", err)
	}

	return &Resolver{
		excludeDirs: excludeDirsMap,
		buildTags:   buildTags,
		log:         log,
		wd:          wd,
	}, nil
}

func (r Resolver) isIgnoredDir(dir string) bool {
	cleanName := filepath.Clean(dir)

	dirName := filepath.Base(cleanName)

	// https://github.com/golang/dep/issues/298
	// https://github.com/tools/godep/issues/140
	if strings.HasPrefix(dirName, ".") && dirName != "." && dirName != ".." {
		return true
	}
	if strings.HasPrefix(dirName, "_") {
		return true
	}

	for _, dirExludeRe := range r.excludeDirs {
		if dirExludeRe.MatchString(cleanName) {
			return true
		}
	}

	return false
}

func (r *Resolver) resolveRecursively(root string, prog *Program) error {
	// import root
	if err := r.resolveDir(root, prog); err != nil {
		return err
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("can't read dir %s: %s", root, err)
	}
	// TODO: pass cached fis to build.Context

	for _, fi := range fis {
		if !fi.IsDir() {
			// ignore files: they were already imported by resolveDir(root)
			continue
		}

		subdir := filepath.Join(root, fi.Name())

		// Normalize each subdir because working directory can be one of these subdirs:
		// working dir = /app/subdir, resolve root is ../, without this normalization
		// path of subdir will be "../subdir" but it must be ".".
		// Normalize path before checking is ignored dir.
		subdir, err := r.normalizePath(subdir)
		if err != nil {
			return err
		}

		if r.isIgnoredDir(subdir) {
			r.skippedDirs = append(r.skippedDirs, subdir)
			continue
		}

		if err := r.resolveRecursively(subdir, prog); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveDir(dir string, prog *Program) error {
	// TODO: fork build.Import to reuse AST parsing
	bp, err := prog.bctx.ImportDir(dir, build.ImportComment|build.IgnoreVendor)
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return nil
		}

		err = fmt.Errorf("can't import dir %q: %s", dir, err)
		r.importErrorsOccured++
		if r.importErrorsOccured >= 10 {
			return err
		}

		r.log.Warnf("Can't analyze dir %q: %s", dir, err)
		return nil
	}

	pkg := Package{
		bp: bp,
	}
	prog.addPackage(&pkg)
	return nil
}

func (r Resolver) addFakePackage(filePath string, prog *Program) {
	// Don't take build tags, is it test file or not, etc
	// into account. If user explicitly wants to analyze this file
	// do it.
	p := Package{
		bp: &build.Package{
			// TODO: detect is it test file or not: without that we can't analyze only one test file
			GoFiles: []string{filePath},
		},
		isFake: true,
		dir:    filepath.Dir(filePath),
	}
	prog.addPackage(&p)
}

func (r Resolver) Resolve(paths ...string) (prog *Program, err error) {
	startedAt := time.Now()
	defer func() {
		r.log.Infof("Paths resolving took %s: %s", time.Since(startedAt), prog)
	}()

	if len(paths) == 0 {
		return nil, fmt.Errorf("no paths are set")
	}

	bctx := build.Default
	bctx.BuildTags = append(bctx.BuildTags, r.buildTags...)
	prog = &Program{
		bctx: bctx,
	}

	for _, path := range paths {
		if err := r.resolvePath(path, prog); err != nil {
			return nil, err
		}
	}

	if len(r.skippedDirs) != 0 {
		r.log.Infof("Skipped dirs: %s", r.skippedDirs)
	}

	return prog, nil
}

func (r *Resolver) normalizePath(path string) (string, error) {
	return fsutils.ShortestRelPath(path, r.wd)
}

func (r *Resolver) resolvePath(path string, prog *Program) error {
	needRecursive := strings.HasSuffix(path, "/...")
	if needRecursive {
		path = filepath.Dir(path)
	}

	evalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("can't eval symlinks for path %s: %s", path, err)
	}
	path = evalPath

	path, err = r.normalizePath(path)
	if err != nil {
		return err
	}

	if needRecursive {
		if err = r.resolveRecursively(path, prog); err != nil {
			return fmt.Errorf("can't recursively resolve %s: %s", path, err)
		}

		return nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("can't find path %s: %s", path, err)
	}

	if fi.IsDir() {
		if err := r.resolveDir(path, prog); err != nil {
			return fmt.Errorf("can't resolve dir %s: %s", path, err)
		}
		return nil
	}

	r.addFakePackage(path, prog)
	return nil
}
