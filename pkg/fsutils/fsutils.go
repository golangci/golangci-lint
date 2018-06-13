package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func ShortestRelPath(path string, wd string) (string, error) {
	if wd == "" { // get it if user don't have cached working dir
		var err error
		wd, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("can't get working directory: %s", err)
		}
	}

	// make path absolute and then relative to be able to fix this case:
	// we'are in /test dir, we want to normalize ../test, and have file file.go in this dir;
	// it must have normalized path file.go, not ../test/file.go,
	var absPath string
	if filepath.IsAbs(path) {
		absPath = path
	} else {
		absPath = filepath.Join(wd, path)
	}

	relPath, err := filepath.Rel(wd, absPath)
	if err != nil {
		return "", fmt.Errorf("can't get relative path for path %s and root %s: %s",
			absPath, wd, err)
	}

	return relPath, nil
}
