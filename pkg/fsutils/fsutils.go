package fsutils

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func GetProjectRoot() string {
	return path.Join(build.Default.GOPATH, "src", "github.com", "golangci", "golangci-worker")
}

type ProjectPaths struct {
	Files     []string
	Dirs      []string
	IsDirsRun bool
}

func (p ProjectPaths) MixedPaths() []string {
	if p.IsDirsRun {
		return p.Dirs
	}

	return p.Files
}

func (p ProjectPaths) FilesGrouppedByDirs() [][]string {
	dirToFiles := map[string][]string{}
	for _, f := range p.Files {
		dir := filepath.Dir(f)
		dirToFiles[dir] = append(dirToFiles[dir], f)
	}

	ret := [][]string{}
	for _, files := range dirToFiles {
		ret = append(ret, files)
	}
	return ret
}

func processPaths(root string, paths []string, maxPaths int) ([]string, error) {
	if len(paths) > maxPaths {
		logrus.Warnf("Gofmt: got too much paths (%d), analyze first %d", len(paths), maxPaths)
		paths = paths[:maxPaths]
	}

	ret := []string{}
	for _, p := range paths {
		if !filepath.IsAbs(p) {
			ret = append(ret, p)
			continue
		}

		relPath, err := filepath.Rel(root, p)
		if err != nil {
			return nil, fmt.Errorf("can't get relative path for path %s and root %s: %s",
				p, root, err)
		}
		ret = append(ret, relPath)
	}

	return ret, nil
}

func GetPathsForAnalysis(ctx context.Context, inputPaths []string, includeTests bool) (ret *ProjectPaths, err error) {
	defer func(startedAt time.Time) {
		if ret != nil {
			logrus.Infof("Found paths for analysis for %s: %s", time.Since(startedAt), ret.MixedPaths())
		}
	}(time.Now())

	for _, path := range inputPaths {
		if strings.HasSuffix(path, ".go") && len(inputPaths) != 1 {
			return nil, fmt.Errorf("Specific files for analysis are allowed only if one file is set")
		}
	}

	excludeDirs := []string{"vendor", "testdata", "examples", "Godeps", "builtin"}
	pr := NewPathResolver(excludeDirs, []string{".go"}, includeTests)
	paths, err := pr.Resolve(inputPaths...)
	if err != nil {
		return nil, fmt.Errorf("can't resolve paths: %s", err)
	}

	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working dir: %s", err)
	}

	files, err := processPaths(root, paths.Files(), 10000)
	if err != nil {
		return nil, fmt.Errorf("can't process resolved files: %s", err)
	}

	dirs, err := processPaths(root, paths.Dirs(), 1000)
	if err != nil {
		return nil, fmt.Errorf("can't process resolved dirs: %s", err)
	}

	for i := range dirs {
		dir := dirs[i]
		if dir != "." && !filepath.IsAbs(dir) {
			dirs[i] = "./" + dir
		}
	}

	return &ProjectPaths{
		Files:     files,
		Dirs:      dirs,
		IsDirsRun: len(dirs) != 0,
	}, nil
}

func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}
