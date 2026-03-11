package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/timeutils"
)

func setupCache(t *testing.T) *Cache {
	t.Helper()

	log := logutils.NewStderrLog("skip")
	sw := timeutils.NewStopwatch("pkgcache", log)

	pkgCache, err := NewCache(sw, log)
	require.NoError(t, err)

	return pkgCache
}

func fakePackage() *packages.Package {
	return &packages.Package{
		PkgPath: "github.com/golangci/example",
		Dir:     "./testdata",
		Module: &packages.Module{
			Path: "github.com/golangci/example",
			Dir:  ".",
		},
		CompiledGoFiles: []string{
			"./testdata/hello.go",
		},
		Imports: map[string]*packages.Package{
			"a": {
				PkgPath: "github.com/golangci/example/a",
			},
			"b": {
				PkgPath: "github.com/golangci/example/b",
			},
			"unsafe": {
				PkgPath: "unsafe",
			},
		},
	}
}

type Foo struct {
	Value string
}

func TestCache_Put(t *testing.T) {
	t.Setenv("GOLANGCI_LINT_CACHE", t.TempDir())

	pkgCache := setupCache(t)

	pkg := fakePackage()

	in := &Foo{Value: "hello"}

	err := pkgCache.Put(pkg, HashModeNeedAllDeps, "key", in)
	require.NoError(t, err)

	out := &Foo{}
	err = pkgCache.Get(pkg, HashModeNeedAllDeps, "key", out)
	require.NoError(t, err)

	assert.Equal(t, in, out)

	pkgCache.Close()
}

func TestCache_Get_missing_data(t *testing.T) {
	t.Setenv("GOLANGCI_LINT_CACHE", t.TempDir())

	pkgCache := setupCache(t)

	pkg := fakePackage()

	out := &Foo{}
	err := pkgCache.Get(pkg, HashModeNeedAllDeps, "key", out)
	require.Error(t, err)

	require.ErrorIs(t, err, ErrMissing)

	pkgCache.Close()
}

func TestCache_buildKey(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	actionID, err := pkgCache.buildKey(pkg, HashModeNeedAllDeps, "")
	require.NoError(t, err)

	assert.Equal(t, "f9cad919ea6b70342f66a4bfa689cc1c14bca2cefea26aa1bcb6934ded8b86a8", fmt.Sprintf("%x", actionID))
}

func TestCache_pkgActionID(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	actionID, err := pkgCache.pkgActionID(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "408e9acfe5f8a37af1c4d244a44f1fa416e85397002ce10d44dffd9bb24c1db5", fmt.Sprintf("%x", actionID))
}

func TestCache_packageHash_load(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	pkgCache.pkgHashes.Store(pkg, hashResults{HashModeNeedAllDeps: "fake"})

	hash, err := pkgCache.packageHash(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "fake", hash)
}

func TestCache_packageHash_store(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	hash, err := pkgCache.packageHash(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "f4825f86811c1ee286edbd347d556bdcb9e4fbddeb6fdd61633ab61b80a1bfdf", hash)

	results, ok := pkgCache.pkgHashes.Load(pkg)
	require.True(t, ok)

	hashRes := results.(hashResults)

	require.Len(t, hashRes, 3)

	assert.Equal(t, "1ee7a6dda5a5ab959e893844bfb1e456daca72f55c38f900b82e9324cfc84eb9", hashRes[HashModeNeedOnlySelf])
	assert.Equal(t, "6b7d112bb0bd2834cbc7c3c58ab7bf580bf51c5e0fb5fb366caf4f3d189aded6", hashRes[HashModeNeedDirectDeps])
	assert.Equal(t, "f4825f86811c1ee286edbd347d556bdcb9e4fbddeb6fdd61633ab61b80a1bfdf", hashRes[HashModeNeedAllDeps])
}

func TestCache_computeHash(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	results, err := pkgCache.computePkgHash(pkg)
	require.NoError(t, err)

	require.Len(t, results, 3)

	assert.Equal(t, "1ee7a6dda5a5ab959e893844bfb1e456daca72f55c38f900b82e9324cfc84eb9", results[HashModeNeedOnlySelf])
	assert.Equal(t, "6b7d112bb0bd2834cbc7c3c58ab7bf580bf51c5e0fb5fb366caf4f3d189aded6", results[HashModeNeedDirectDeps])
	assert.Equal(t, "f4825f86811c1ee286edbd347d556bdcb9e4fbddeb6fdd61633ab61b80a1bfdf", results[HashModeNeedAllDeps])
}

func TestCache_computeHash_samePackageAcrossWorktrees(t *testing.T) {
	t.Setenv("GOLANGCI_LINT_CACHE", t.TempDir())

	rootA := filepath.Join(t.TempDir(), "repo-a")
	rootB := filepath.Join(t.TempDir(), "repo-b")

	pkgA, depA := createPackageFixture(t, rootA)
	pkgB, depB := createPackageFixture(t, rootB)

	pkgCache := setupCache(t)

	rootHashesA, err := pkgCache.computePkgHash(pkgA)
	require.NoError(t, err)

	rootHashesB, err := pkgCache.computePkgHash(pkgB)
	require.NoError(t, err)

	depHashesA, err := pkgCache.computePkgHash(depA)
	require.NoError(t, err)

	depHashesB, err := pkgCache.computePkgHash(depB)
	require.NoError(t, err)

	assert.Equal(t, depHashesA, depHashesB)
	assert.Equal(t, rootHashesA, rootHashesB)
}

func createPackageFixture(t *testing.T, root string) (pkg, dep *packages.Package) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Join(root, "dep"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "pkg"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "dep", "dep.go"), []byte("package dep\n\nconst Name = \"dep\"\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "pkg", "main.go"), []byte("package pkg\n\nimport _ \"example.com/project/dep\"\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "pkg", "extra.go"), []byte("package pkg\n\nconst value = 1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "pkg", "ignored.go"), []byte("//go:build ignored\n\npackage pkg\n"), 0o600))

	module := &packages.Module{
		Path: "example.com/project",
		Dir:  root,
	}

	dep = &packages.Package{
		PkgPath: "example.com/project/dep",
		Dir:     filepath.Join(root, "dep"),
		Module:  module,
		CompiledGoFiles: []string{
			filepath.Join(root, "dep", "dep.go"),
		},
	}

	pkg = &packages.Package{
		PkgPath: "example.com/project/pkg",
		Dir:     filepath.Join(root, "pkg"),
		Module:  module,
		CompiledGoFiles: []string{
			filepath.Join(root, "pkg", "main.go"),
			filepath.Join(root, "pkg", "extra.go"),
		},
		IgnoredFiles: []string{
			filepath.Join(root, "pkg", "ignored.go"),
		},
		Imports: map[string]*packages.Package{
			"example.com/project/dep": dep,
		},
	}

	return pkg, dep
}
