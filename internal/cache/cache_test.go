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

	assert.Equal(t, "f32bf1bf010aa9b570e081c64ec9e22e17aafa1e822990ba952905ec5fdf8d9d", fmt.Sprintf("%x", actionID))
}

func TestCache_pkgActionID(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	actionID, err := pkgCache.pkgActionID(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "f690f05acd1024386ae912d9ad9c04080523b9a899f6afe56ab3108d88215c1d", fmt.Sprintf("%x", actionID))
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

	assert.Equal(t, "9c602ef861197b6807e82c99caa7c4042eb03c1a92886303fb02893744355131", hash)

	results, ok := pkgCache.pkgHashes.Load(pkg)
	require.True(t, ok)

	hashRes := results.(hashResults)

	require.Len(t, hashRes, 3)

	assert.Equal(t, "8978e3d76c6f99e9663558d7147a7790f229a676804d1fde706a611898547b74", hashRes[HashModeNeedOnlySelf])
	assert.Equal(t, "b1aef902a0619b5cbfc2d6e2e91a73dd58dd448e58274b2d7a5ff8efd97aefa4", hashRes[HashModeNeedDirectDeps])
	assert.Equal(t, "9c602ef861197b6807e82c99caa7c4042eb03c1a92886303fb02893744355131", hashRes[HashModeNeedAllDeps])
}

func TestCache_computeHash(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	results, err := pkgCache.computePkgHash(pkg)
	require.NoError(t, err)

	require.Len(t, results, 3)

	assert.Equal(t, "8978e3d76c6f99e9663558d7147a7790f229a676804d1fde706a611898547b74", results[HashModeNeedOnlySelf])
	assert.Equal(t, "b1aef902a0619b5cbfc2d6e2e91a73dd58dd448e58274b2d7a5ff8efd97aefa4", results[HashModeNeedDirectDeps])
	assert.Equal(t, "9c602ef861197b6807e82c99caa7c4042eb03c1a92886303fb02893744355131", results[HashModeNeedAllDeps])
}

func TestCache_computeHash_module(t *testing.T) {
	pkgCache := setupCache(t)

	// This creates a new random directory for each test run, but the hash will stay consistent.
	tempDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tempDir, "foo.go"), []byte("package foo"), 0600)
	require.NoError(t, err)

	pkg := fakePackage()
	pkg.Module = &packages.Module{
		Dir:  tempDir,
		Path: "github.com/golangci/example",
	}
	pkg.IgnoredFiles = []string{filepath.Join(tempDir, "foo.go")}

	results, err := pkgCache.computePkgHash(pkg)
	require.NoError(t, err)

	require.Len(t, results, 3)

	assert.Equal(t, "ac5d79f4630d6b5f1e4ac88bfa9974698ea44ea1f53760bc75001bd7c9ce9064", results[HashModeNeedOnlySelf])
	assert.Equal(t, "73fd7ef46b20efdd1bff7eaedade26b1655ee727d83f4210b7934bf4d4a1ac1d", results[HashModeNeedDirectDeps])
	assert.Equal(t, "063ffbaa2ea6b2a9a02724effe2a616b88672a4d8c3c4c6c8d777bfe5a49e7b4", results[HashModeNeedAllDeps])
}
