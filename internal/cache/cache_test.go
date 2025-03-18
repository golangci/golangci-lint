package cache

import (
	"fmt"
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

	assert.Equal(t, "4db1e9671b1800244a938e1aabe6db77495a5dff82acb434ab435eaabff59372", fmt.Sprintf("%x", actionID))
}

func TestCache_pkgActionID(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	actionID, err := pkgCache.pkgActionID(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "6ab6e6fea09b9390f266880ab504e267ea12b76b14a783013843effb1d7d31a5", fmt.Sprintf("%x", actionID))
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

	assert.Equal(t, "11d5b95830a3d86ad6dcac15b0adb46f686a30f1633cc1b74f7133d093da7979", hash)

	results, ok := pkgCache.pkgHashes.Load(pkg)
	require.True(t, ok)

	hashRes := results.(hashResults)

	require.Len(t, hashRes, 3)

	assert.Equal(t, "855344dcd4884d73c77e2eede91270d871e18b57aeb170007b5c682b39e52c96", hashRes[HashModeNeedOnlySelf])
	assert.Equal(t, "f1288ac9e75477e1ed42c06ba7d78213ed3cc52799c01a4f695bc54cd4e8a698", hashRes[HashModeNeedDirectDeps])
	assert.Equal(t, "11d5b95830a3d86ad6dcac15b0adb46f686a30f1633cc1b74f7133d093da7979", hashRes[HashModeNeedAllDeps])
}

func TestCache_computeHash(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	results, err := pkgCache.computePkgHash(pkg)
	require.NoError(t, err)

	require.Len(t, results, 3)

	assert.Equal(t, "855344dcd4884d73c77e2eede91270d871e18b57aeb170007b5c682b39e52c96", results[HashModeNeedOnlySelf])
	assert.Equal(t, "f1288ac9e75477e1ed42c06ba7d78213ed3cc52799c01a4f695bc54cd4e8a698", results[HashModeNeedDirectDeps])
	assert.Equal(t, "11d5b95830a3d86ad6dcac15b0adb46f686a30f1633cc1b74f7133d093da7979", results[HashModeNeedAllDeps])
}
