// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"fmt"
	base "log"
	"os"
	"path/filepath"
	"sync"
)

// Default returns the default cache to use.
// It never returns nil.
func Default() Cache {
	return initDefaultCacheOnce()
}

var initDefaultCacheOnce = sync.OnceValue(initDefaultCache)

// cacheREADME is a message stored in a README in the cache directory.
// Because the cache lives outside the normal Go trees, we leave the
// README as a courtesy to explain where it came from.
const cacheREADME = `This directory holds cached build artifacts from golangci-lint.
`

// initDefaultCache does the work of finding the default cache
// the first time Default is called.
func initDefaultCache() Cache {
	dir, _, err := DefaultDir()
	if err != nil {
		base.Fatalf("build cache is required, but could not be located: %v", err)
	}
	if dir == "off" {
		if defaultDirErr != nil {
			base.Fatalf("build cache is required, but could not be located: %v", defaultDirErr)
		}
		base.Fatalf("build cache is disabled by %s=off, but required as of Go 1.12", envGolangciLintCache)
	}
	if err := os.MkdirAll(dir, 0o777); err != nil {
		base.Fatalf("failed to initialize build cache at %s: %s\n", dir, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "README")); err != nil {
		// Best effort.
		os.WriteFile(filepath.Join(dir, "README"), []byte(cacheREADME), 0666)
	}

	diskCache, err := Open(dir)
	if err != nil {
		base.Fatalf("failed to initialize build cache at %s: %s\n", dir, err)
	}

	if v := os.Getenv(envGolangciLintCacheProg); v != "" {
		return startCacheProg(v, diskCache)
	}

	return diskCache
}

var (
	defaultDirOnce    sync.Once
	defaultDir        string
	defaultDirChanged bool // effective value differs from $GOLANGCI_LINT_CACHE
	defaultDirErr     error
)

// DefaultDir returns the effective GOLANGCI_LINT_CACHE setting.
// It returns "off" if the cache is disabled,
// and reports whether the effective value differs from GOLANGCI_LINT_CACHE.
func DefaultDir() (string, bool, error) {
	// Save the result of the first call to DefaultDir for later use in
	// initDefaultCache. cmd/go/main.go explicitly sets GOLANGCI_LINT_CACHE so that
	// subprocesses will inherit it, but that means initDefaultCache can't
	// otherwise distinguish between an explicit "off" and a UserCacheDir error.

	defaultDirOnce.Do(func() {
		// Compute default location.
		dir, err := os.UserCacheDir()
		if err != nil {
			defaultDir = "off"
			defaultDirErr = fmt.Errorf("%s is not defined and %v", envGolangciLintCache, err)
		} else {
			defaultDir = filepath.Join(dir, "golangci-lint")
		}

		newDir := os.Getenv(envGolangciLintCache)
		if newDir != "" {
			defaultDirErr = nil
			defaultDirChanged = newDir != defaultDir
			defaultDir = newDir
			if filepath.IsAbs(defaultDir) || defaultDir == "off" {
				return
			}
			defaultDir = "off"
			defaultDirErr = fmt.Errorf("%s is not an absolute path", envGolangciLintCache)
			return
		}
	})

	return defaultDir, defaultDirChanged, defaultDirErr
}
