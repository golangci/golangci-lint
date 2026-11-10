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

	assert.Equal(t, "4b048b43c3506df16209043a08cdd4498b471590cddd14d89de39907a87e9d1d", fmt.Sprintf("%x", actionID))
}

func TestCache_pkgActionID(t *testing.T) {
	pkgCache := setupCache(t)

	pkg := fakePackage()

	actionID, err := pkgCache.pkgActionID(pkg, HashModeNeedAllDeps)
	require.NoError(t, err)

	assert.Equal(t, "cd892bf2556f11d7eb91626f5e228859646a001a9bd489f3d267ef687879ccf2", fmt.Sprintf("%x", actionID))
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

	assert.Equal(t, "9114847c7d8ad5ee2e3f04f989a069fcdd91f762f220b8b16cad72de68df3530", hash)

	results, ok := pkgCache.pkgHashes.Load(pkg)
	require.True(t, ok)

	hashRes := results.(hashResults)

	require.Len(t, hashRes, 3)

	assert.Equal(t, "3699dbf7b96814ac7f486053ce9de0409a360f3b6990a417e656fd7452043dc3", hashRes[HashModeNeedOnlySelf])
	assert.Equal(t, "3a1961cdd38edf42038087ec5e6b44176c8d40b8c04598ef4599f38e21801ad6", hashRes[HashModeNeedDirectDeps])
	assert.Equal(t, "9114847c7d8ad5ee2e3f04f989a069fcdd91f762f220b8b16cad72de68df3530", hashRes[HashModeNeedAllDeps])
}

func TestCache_computeHash(t *testing.T) {
	pkgCache := setupCache(t)

	testCases := []struct {
		desc string
		pkg  func() *packages.Package

		expectedHashModeNeedOnlySelf   string
		expectedHashModeNeedDirectDeps string
		expectedHashModeNeedAllDeps    string
	}{
		{
			desc: "case a",
			pkg: func() *packages.Package {
				pkg := fakePackage()

				pkg.CompiledGoFiles = []string{
					"./testdata/hello.go",
					"./testdata/world.go",
				}

				pkg.IgnoredFiles = []string{}

				return pkg
			},
			expectedHashModeNeedOnlySelf:   "3aa2ecfb2d98dbc96766d158663942aa72ff354e388b9df61f008726fbf8af94",
			expectedHashModeNeedDirectDeps: "542a5cbc0c2d710b35a5ea07dadaa3d562010e7d191a91d0e0eb66421df17d05",
			expectedHashModeNeedAllDeps:    "5d7d365ffb255000579a39d580ce0674d1bef7f2d3a4b70da01c92e4795ef0f8",
		},
		{
			desc: "case b",
			pkg: func() *packages.Package {
				pkg := fakePackage()

				pkg.CompiledGoFiles = []string{
					"./testdata/hello.go",
				}

				pkg.IgnoredFiles = []string{
					"./testdata/world.go",
				}

				return pkg
			},
			expectedHashModeNeedOnlySelf:   "ad81fd8c68683d7445ee55c8284041e4ad2f42be429ed35d0583f0da69c7f1ae",
			expectedHashModeNeedDirectDeps: "48b3430ddb0300030cdd30d57390afde734adafe3e4d09b70f97614f9826c5e5",
			expectedHashModeNeedAllDeps:    "ffa6661e52a4fcdab4c1c378677f8ff9b89cdb92712cf0d7fb3658df7147d3de",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			results, err := pkgCache.computePkgHash(test.pkg())
			require.NoError(t, err)

			require.Len(t, results, 3)

			assert.Equal(t, test.expectedHashModeNeedOnlySelf, results[HashModeNeedOnlySelf])
			assert.Equal(t, test.expectedHashModeNeedDirectDeps, results[HashModeNeedDirectDeps])
			assert.Equal(t, test.expectedHashModeNeedAllDeps, results[HashModeNeedAllDeps])
		})
	}
}
