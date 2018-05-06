package fsutils

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-shared/pkg/analytics"
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

func processPaths(root string, paths []string, maxPaths int) ([]string, error) {
	if len(paths) > maxPaths {
		analytics.Log(context.TODO()).Warnf("Gofmt: got too much paths (%d), analyze first %d", len(paths), maxPaths)
		paths = paths[:maxPaths]
	}

	ret := []string{}
	for i := range paths {
		if !filepath.IsAbs(paths[i]) {
			ret = append(ret, paths[i])
			continue
		}

		relPath, err := filepath.Rel(root, paths[i])
		if err != nil {
			return nil, fmt.Errorf("can't get relative path for path %s and root %s: %s",
				paths[i], root, err)
		}
		ret = append(ret, relPath)
	}

	return ret, nil
}

func GetPathsForAnalysis(inputPaths []string) (*ProjectPaths, error) {
	for _, path := range inputPaths {
		if strings.HasSuffix(path, ".go") && len(inputPaths) != 1 {
			return nil, fmt.Errorf("Specific files for analysis are allowed only if one file is set")
		}
	}

	excludeDirs := []string{"vendor", "testdata", "examples", "Godeps"}
	pr := NewPathResolver(excludeDirs, []string{".go"})
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
		if dir != "." {
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
