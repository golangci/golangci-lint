package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/internal/go/cache"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/timeutils"
)

type HashMode int

const (
	HashModeNeedOnlySelf HashMode = iota
	HashModeNeedDirectDeps
	HashModeNeedAllDeps
)

var ErrMissing = errors.New("missing data")

type hashResults map[HashMode]string

// Cache is a per-package data cache.
// A cached data is invalidated when package,
// or it's dependencies change.
type Cache struct {
	lowLevelCache cache.Cache
	pkgHashes     sync.Map
	sw            *timeutils.Stopwatch
	log           logutils.Log
	ioSem         chan struct{} // semaphore limiting parallel IO
}

func NewCache(sw *timeutils.Stopwatch, log logutils.Log) (*Cache, error) {
	return &Cache{
		lowLevelCache: cache.Default(),
		sw:            sw,
		log:           log,
		ioSem:         make(chan struct{}, runtime.GOMAXPROCS(-1)),
	}, nil
}

func (c *Cache) Close() {
	err := c.sw.TrackStageErr("close", c.lowLevelCache.Close)
	if err != nil {
		c.log.Errorf("cache close: %v", err)
	}
}

// PrecomputePackageHashes walks the package dependency graph in topological
// order (leaves first) and computes hashes for all packages concurrently.
// This warms the pkgHashes cache so that subsequent Get/Put calls don't
// need to recursively compute hashes on the critical path.
func (c *Cache) PrecomputePackageHashes(pkgs []*packages.Package) {
	// Collect all packages in the dependency graph.
	all := map[*packages.Package]struct{}{}
	var collect func(pkg *packages.Package)
	collect = func(pkg *packages.Package) {
		if _, seen := all[pkg]; seen {
			return
		}
		all[pkg] = struct{}{}
		for _, dep := range pkg.Imports {
			collect(dep)
		}
	}
	for _, pkg := range pkgs {
		collect(pkg)
	}

	if len(all) == 0 {
		return
	}

	// Build reverse dep edges and count pending deps per package.
	type pkgState struct {
		pending atomic.Int32
	}
	states := make(map[*packages.Package]*pkgState, len(all))
	reverseDeps := make(map[*packages.Package][]*packages.Package, len(all))

	for pkg := range all {
		states[pkg] = &pkgState{}
	}

	for pkg := range all {
		for _, dep := range pkg.Imports {
			if dep.PkgPath == "unsafe" {
				continue
			}
			if _, ok := all[dep]; ok {
				states[pkg].pending.Add(1)
				reverseDeps[dep] = append(reverseDeps[dep], pkg)
			}
		}
	}

	// Track overall completion.
	var remaining sync.WaitGroup
	remaining.Add(len(all))

	ch := make(chan *packages.Package, len(all))

	// Seed with leaf packages.
	for pkg, st := range states {
		if st.pending.Load() == 0 {
			ch <- pkg
		}
	}

	// Workers process packages whose deps are all resolved.
	var workers sync.WaitGroup
	for range runtime.GOMAXPROCS(-1) {
		workers.Go(func() {
			for pkg := range ch {
				// Deps are already in pkgHashes, so this won't recurse deeply.
				_, _ = c.packageHash(pkg, HashModeNeedAllDeps)
				remaining.Done()

				for _, rdep := range reverseDeps[pkg] {
					if states[rdep].pending.Add(-1) == 0 {
						ch <- rdep
					}
				}
			}
		})
	}

	remaining.Wait()
	close(ch)
	workers.Wait()
}

func (c *Cache) Put(pkg *packages.Package, mode HashMode, key string, data any) error {
	buf, err := c.encode(data)
	if err != nil {
		return err
	}

	actionID, err := c.buildKey(pkg, mode, key)
	if err != nil {
		return fmt.Errorf("failed to calculate package %s action id: %w", pkg.Name, err)
	}

	err = c.putBytes(actionID, buf)
	if err != nil {
		return fmt.Errorf("failed to save data to low-level cache by key %s for package %s: %w", key, pkg.Name, err)
	}

	return nil
}

func (c *Cache) Get(pkg *packages.Package, mode HashMode, key string, data any) error {
	actionID, err := c.buildKey(pkg, mode, key)
	if err != nil {
		return fmt.Errorf("failed to calculate package %s action id: %w", pkg.Name, err)
	}

	cachedData, err := c.getBytes(actionID)
	if err != nil {
		if cache.IsErrMissing(err) {
			return ErrMissing
		}
		return fmt.Errorf("failed to get data from low-level cache by key %s for package %s: %w", key, pkg.Name, err)
	}

	return c.decode(cachedData, data)
}

func (c *Cache) buildKey(pkg *packages.Package, mode HashMode, key string) (cache.ActionID, error) {
	return timeutils.TrackStage(c.sw, "key build", func() (cache.ActionID, error) {
		actionID, err := c.pkgActionID(pkg, mode)
		if err != nil {
			return actionID, err
		}

		subkey, subkeyErr := cache.Subkey(actionID, key)
		if subkeyErr != nil {
			return actionID, fmt.Errorf("failed to build subkey: %w", subkeyErr)
		}

		return subkey, nil
	})
}

func (c *Cache) pkgActionID(pkg *packages.Package, mode HashMode) (cache.ActionID, error) {
	hash, err := c.packageHash(pkg, mode)
	if err != nil {
		return cache.ActionID{}, fmt.Errorf("failed to get package hash: %w", err)
	}

	key, err := cache.NewHash("action ID")
	if err != nil {
		return cache.ActionID{}, fmt.Errorf("failed to make a hash: %w", err)
	}

	fmt.Fprintf(key, "pkgpath %s\n", pkg.PkgPath)
	fmt.Fprintf(key, "pkghash %s\n", hash)

	return key.Sum(), nil
}

func (c *Cache) packageHash(pkg *packages.Package, mode HashMode) (string, error) {
	results, found := c.pkgHashes.Load(pkg)
	if found {
		hashRes := results.(hashResults)
		if result, ok := hashRes[mode]; ok {
			return result, nil
		}

		return "", fmt.Errorf("no mode %d in hash result", mode)
	}

	hashRes, err := c.computePkgHash(pkg)
	if err != nil {
		return "", err
	}

	result, found := hashRes[mode]
	if !found {
		return "", fmt.Errorf("invalid mode %d", mode)
	}

	c.pkgHashes.Store(pkg, hashRes)

	return result, nil
}

// computePkgHash computes a package's hash.
// The hash is based on all Go files that make up the package,
// as well as the hashes of imported packages.
func (c *Cache) computePkgHash(pkg *packages.Package) (hashResults, error) {
	key, err := cache.NewHash("package hash")
	if err != nil {
		return nil, fmt.Errorf("failed to make a hash: %w", err)
	}

	hashRes := hashResults{}

	fmt.Fprintf(key, "pkgpath %s\n", pkg.PkgPath)

	// Hash all files in the package concurrently.
	files := slices.Concat(pkg.CompiledGoFiles, pkg.IgnoredFiles)

	type fileHashResult struct {
		name string
		hash [cache.HashSize]byte
	}

	results := make([]fileHashResult, len(files))

	var wg sync.WaitGroup
	var hashErr atomic.Pointer[error]

	for i, f := range files {
		// Pre-compute display name before spawning goroutine.
		name := f
		if pkg.Module != nil && pkg.Module.Version == "" {
			name = pkg.Module.Path + strings.TrimPrefix(filepath.ToSlash(f), filepath.ToSlash(pkg.Module.Dir))
		}
		results[i].name = name

		wg.Go(func() {
			h, fErr := c.fileHash(f)
			if fErr != nil {
				e := fmt.Errorf("failed to calculate file %s hash: %w", f, fErr)
				hashErr.CompareAndSwap(nil, &e)
				return
			}
			results[i].hash = h
		})
	}

	wg.Wait()

	if ep := hashErr.Load(); ep != nil {
		return nil, *ep
	}

	for _, r := range results {
		fmt.Fprintf(key, "file %s %x\n", r.name, r.hash)
	}

	curSum := key.Sum()
	hashRes[HashModeNeedOnlySelf] = hex.EncodeToString(curSum[:])

	imps := slices.SortedFunc(maps.Values(pkg.Imports), func(a, b *packages.Package) int {
		return strings.Compare(a.PkgPath, b.PkgPath)
	})

	if err := c.computeDepsHash(HashModeNeedOnlySelf, imps, key); err != nil {
		return nil, err
	}

	curSum = key.Sum()
	hashRes[HashModeNeedDirectDeps] = hex.EncodeToString(curSum[:])

	if err := c.computeDepsHash(HashModeNeedAllDeps, imps, key); err != nil {
		return nil, err
	}

	curSum = key.Sum()
	hashRes[HashModeNeedAllDeps] = hex.EncodeToString(curSum[:])

	return hashRes, nil
}

func (c *Cache) computeDepsHash(depMode HashMode, imps []*packages.Package, key *cache.Hash) error {
	for _, dep := range imps {
		if dep.PkgPath == "unsafe" {
			continue
		}

		depHash, err := c.packageHash(dep, depMode)
		if err != nil {
			return fmt.Errorf("failed to calculate hash for dependency %s with mode %d: %w", dep.Name, depMode, err)
		}

		fmt.Fprintf(key, "import %s %s\n", dep.PkgPath, depHash)
	}

	return nil
}

func (c *Cache) putBytes(actionID cache.ActionID, buf *bytes.Buffer) error {
	c.ioSem <- struct{}{}

	err := c.sw.TrackStageErr("cache io", func() error {
		return cache.PutBytes(c.lowLevelCache, actionID, buf.Bytes())
	})

	<-c.ioSem

	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) getBytes(actionID cache.ActionID) ([]byte, error) {
	c.ioSem <- struct{}{}

	cachedData, err := timeutils.TrackStage(c.sw, "cache io", func() ([]byte, error) {
		b, _, errGB := cache.GetBytes(c.lowLevelCache, actionID)
		return b, errGB
	})

	<-c.ioSem

	if err != nil {
		return nil, err
	}

	return cachedData, nil
}

func (c *Cache) fileHash(f string) ([cache.HashSize]byte, error) {
	return cache.FileHash(f)
}

func (c *Cache) encode(data any) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	err := c.sw.TrackStageErr("gob", func() error {
		return gob.NewEncoder(buf).Encode(data)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to gob encode: %w", err)
	}

	return buf, nil
}

func (c *Cache) decode(b []byte, data any) error {
	err := c.sw.TrackStageErr("gob", func() error {
		return gob.NewDecoder(bytes.NewReader(b)).Decode(data)
	})
	if err != nil {
		return fmt.Errorf("failed to gob decode: %w", err)
	}

	return nil
}

func SetSalt(b *bytes.Buffer) {
	cache.SetSalt(b.Bytes())
}

func DefaultDir() string {
	cacheDir, _ := cache.DefaultDir()
	return cacheDir
}
