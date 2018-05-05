package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PathResolver struct {
	excludeDirs           map[string]bool
	allowedFileExtensions map[string]bool
}

type pathResolveState struct {
	files map[string]bool
	dirs  map[string]bool
}

func (s *pathResolveState) addFile(path string) {
	s.files[filepath.Clean(path)] = true
}

func (s *pathResolveState) addDir(path string) {
	s.dirs[filepath.Clean(path)] = true
}

type PathResolveResult struct {
	files []string
	dirs  []string
}

func (prr PathResolveResult) Files() []string {
	return prr.files
}

func (prr PathResolveResult) Dirs() []string {
	return prr.dirs
}

func (s pathResolveState) toResult() *PathResolveResult {
	res := &PathResolveResult{
		files: []string{},
		dirs:  []string{},
	}
	for f := range s.files {
		res.files = append(res.files, f)
	}
	for d := range s.dirs {
		res.dirs = append(res.dirs, d)
	}

	sort.Strings(res.files)
	sort.Strings(res.dirs)
	return res
}

func NewPathResolver(excludeDirs, allowedFileExtensions []string) *PathResolver {
	excludeDirsMap := map[string]bool{}
	for _, dir := range excludeDirs {
		excludeDirsMap[dir] = true
	}

	allowedFileExtensionsMap := map[string]bool{}
	for _, fe := range allowedFileExtensions {
		allowedFileExtensionsMap[fe] = true
	}

	return &PathResolver{
		excludeDirs:           excludeDirsMap,
		allowedFileExtensions: allowedFileExtensionsMap,
	}
}

func (pr PathResolver) isIgnoredDir(dir string) bool {
	dirName := filepath.Base(filepath.Clean(dir)) // ignore dirs on any depth level

	// https://github.com/golang/dep/issues/298
	// https://github.com/tools/godep/issues/140
	if strings.HasPrefix(dirName, ".") && dirName != "." {
		return true
	}
	if strings.HasPrefix(dirName, "_") {
		return true
	}

	return pr.excludeDirs[dirName]
}

func (pr PathResolver) isAllowedFile(path string) bool {
	return pr.allowedFileExtensions[filepath.Ext(path)]
}

func (pr PathResolver) resolveRecursively(root string, state *pathResolveState) error {
	walkErr := filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			if pr.isIgnoredDir(p) {
				return filepath.SkipDir
			}
			state.addDir(p)
			return nil
		}

		if pr.isAllowedFile(p) {
			state.addFile(p)
		}
		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("can't walk dir %s: %s", root, walkErr)
	}

	return nil
}

func (pr PathResolver) resolveDir(root string, state *pathResolveState) error {
	walkErr := filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i.IsDir() {
			if filepath.Clean(p) != filepath.Clean(root) {
				return filepath.SkipDir
			}
			state.addDir(p)
			return nil
		}

		if pr.isAllowedFile(p) {
			state.addFile(p)
		}
		return nil
	})

	if walkErr != nil {
		return fmt.Errorf("can't walk dir %s: %s", root, walkErr)
	}

	return nil
}

func (pr PathResolver) Resolve(paths ...string) (*PathResolveResult, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no paths are set")
	}

	state := &pathResolveState{
		files: map[string]bool{},
		dirs:  map[string]bool{},
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "/...") {
			if err := pr.resolveRecursively(filepath.Dir(path), state); err != nil {
				return nil, fmt.Errorf("can't recursively resolve %s: %s", path, err)
			}
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("can't find path %s: %s", path, err)
		}

		if fi.IsDir() {
			if err := pr.resolveDir(path, state); err != nil {
				return nil, fmt.Errorf("can't resolve dir %s: %s", path, err)
			}
			continue
		}

		state.addFile(path)
	}

	return state.toResult(), nil
}
