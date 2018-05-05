package golinters

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-shared/pkg/analytics"
)

type ProjectPaths struct {
	files []string
	dirs  []string
}

func processPaths(root string, paths []string, maxPaths int) ([]string, error) {
	if len(paths) >= maxPaths {
		analytics.Log(context.TODO()).Warnf("Gofmt: got too much paths (%d), analyze first %d", len(paths), maxPaths)
		paths = paths[:maxPaths]
	}

	ret := []string{}
	for i := range paths {
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
	log.Printf("root is %q, paths are %q", root, path.Join(root, "..."))
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
