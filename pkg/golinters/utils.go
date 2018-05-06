package golinters

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-shared/pkg/analytics"
)

type ProjectPaths struct {
	files []string
	dirs  []string
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

func getPathsForGoProject(root string) (*ProjectPaths, error) {
	excludeDirs := []string{"vendor", "testdata", "examples", "Godeps"}
	pr := fsutils.NewPathResolver(excludeDirs, []string{".go"})
	paths, err := pr.Resolve(path.Join(root, "..."))
	if err != nil {
		return nil, fmt.Errorf("can't resolve paths: %s", err)
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
		files: files,
		dirs:  dirs,
	}, nil
}

func GetPathsForAnalysis(inputPaths []string) ([]string, error) {
	excludeDirs := []string{"vendor", "testdata", "examples", "Godeps"}
	pr := fsutils.NewPathResolver(excludeDirs, []string{".go"})
	paths, err := pr.Resolve(inputPaths...)
	if err != nil {
		return nil, fmt.Errorf("can't resolve paths: %s", err)
	}

	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working dir: %s", err)
	}

	if len(paths.Files()) == 1 {
		// special case: only one files was set for analysis
		return processPaths(root, paths.Files(), 1)
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

	return dirs, nil
}

func formatCode(code string, cfg *config.Run) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}
