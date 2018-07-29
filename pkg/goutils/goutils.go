package goutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golangci/golangci-lint/pkg/fsutils"
)

var discoverGoRootOnce sync.Once
var discoveredGoRoot string
var discoveredGoRootError error

func DiscoverGoRoot() (string, error) {
	discoverGoRootOnce.Do(func() {
		discoveredGoRoot, discoveredGoRootError = discoverGoRootImpl()
	})

	return discoveredGoRoot, discoveredGoRootError
}

func discoverGoRootImpl() (string, error) {
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

func InGoRoot() (bool, error) {
	goroot, err := DiscoverGoRoot()
	if err != nil {
		return false, err
	}

	wd, err := fsutils.Getwd()
	if err != nil {
		return false, err
	}

	// TODO: strip, then add slashes
	return strings.HasPrefix(wd, goroot), nil
}

func IsCgoFilename(f string) bool {
	return filepath.Base(f) == "C"
}
