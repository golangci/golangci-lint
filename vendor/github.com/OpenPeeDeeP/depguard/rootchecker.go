package depguard

import (
	"go/build"
)

// RootChecker checks if import paths point to root packages.
type RootChecker struct {
	buildCtx *build.Context
	cache    map[string]bool
}

// NewRootChecker creates a new RootChecker instance using the build.Context
// given, or build.Default.
func NewRootChecker(buildCtx *build.Context) *RootChecker {
	// Use the &build.Default if build.Context is not specified
	ctx := buildCtx
	if ctx == nil {
		ctx = &build.Default
	}
	return &RootChecker{
		buildCtx: ctx,
		cache:    make(map[string]bool, 0),
	}
}

// IsRoot checks if the given import path (imported from sourceDir)
// points to a a root package. Subsequent calls with the same arguments
// are cached. This is not thread-safe.
func (rc *RootChecker) IsRoot(path, sourceDir string) (bool, error) {
	key := path + ":::" + sourceDir
	isRoot, ok := rc.cache[key]
	if ok {
		return isRoot, nil
	}
	isRoot, err := rc.calcIsRoot(path, sourceDir)
	if err != nil {
		return false, err
	}
	rc.cache[key] = isRoot
	return isRoot, nil
}

// calcIsRoot performs the call to the build context to check if
// the import path points to a root package.
func (rc *RootChecker) calcIsRoot(path, sourceDir string) (bool, error) {
	pkg, err := rc.buildCtx.Import(path, sourceDir, build.FindOnly)
	if err != nil {
		return false, err
	}
	return pkg.Goroot, nil
}
