package goanalysis

import (
	"errors"
	"fmt"
	"go/types"
	"io"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/objectpath"

	"github.com/golangci/golangci-lint/v2/internal/cache"
)

type Fact struct {
	// PkgPath is the import path of the package owning the fact's object (object facts)
	// or the package the fact is about (package facts).
	// It is empty when the fact belongs to the package being cached,
	// which keeps backward compatibility with the previous cache format.
	PkgPath string
	Path    string // non-empty only for object facts
	Fact    analysis.Fact
}

func (act *action) loadCachedFacts() bool {
	if act.loadCachedFactsDone { // can't be set in parallel
		return act.loadCachedFactsOk
	}

	res := func() bool {
		if act.isInitialPkg {
			return true // load cached facts only for non-initial packages
		}

		if len(act.Analyzer.FactTypes) == 0 {
			return true // no need to load facts
		}

		return act.loadPersistedFacts()
	}()

	act.loadCachedFactsDone = true
	act.loadCachedFactsOk = res

	return res
}

func (act *action) persistFactsToCache() error {
	analyzer := act.Analyzer

	if len(analyzer.FactTypes) == 0 {
		return nil
	}

	// Merge new facts into the package and persist them.
	//
	// We must persist not only the facts about this package's own objects,
	// but also the facts about objects from other packages that are reachable through this package's export data (see exportedFrom).
	// When a package is later restored from the cache instead of being analyzed from source,
	// it re-exports these inherited facts to its own dependents.
	// Dropping them here makes linters using facts silently miss issues for downstream packages analyzed from source.
	// e.g. it makes nolintlint report the related `nolint` directives as unused.

	var facts []Fact

	for key, fact := range act.packageFacts {
		pkgPath := ""
		if key.pkg != act.Package.Types {
			pkgPath = key.pkg.Path()
		}

		facts = append(facts, Fact{
			PkgPath: pkgPath,
			Path:    "",
			Fact:    fact,
		})
	}

	for key, fact := range act.objectFacts {
		obj := key.obj

		// Keep facts about objects reachable downstream through the export data,
		// mirroring the filter used by inheritFacts.
		if !exportedFrom(obj, act.Package.Types) {
			continue
		}

		path, err := objectpath.For(obj)
		if err != nil {
			// The object is not globally addressable
			continue
		}

		pkgPath := ""
		if obj.Pkg() != nil && obj.Pkg() != act.Package.Types {
			pkgPath = obj.Pkg().Path()
		}

		facts = append(facts, Fact{
			PkgPath: pkgPath,
			Path:    string(path),
			Fact:    fact,
		})
	}

	factsCacheDebugf("Caching %d facts for package %q and analyzer %s", len(facts), act.Package.Name, act.Analyzer.Name)

	return act.runner.pkgCache.Put(act.Package, cache.HashModeNeedAllDeps, factCacheKey(analyzer), facts)
}

func (act *action) loadPersistedFacts() bool {
	var facts []Fact

	err := act.runner.pkgCache.Get(act.Package, cache.HashModeNeedAllDeps, factCacheKey(act.Analyzer), &facts)
	if err != nil {
		if !errors.Is(err, cache.ErrMissing) && !errors.Is(err, io.EOF) {
			act.runner.log.Warnf("Failed to get persisted facts: %s", err)
		}

		factsCacheDebugf("No cached facts for package %q and analyzer %s", act.Package.Name, act.Analyzer.Name)

		return false
	}

	factsCacheDebugf("Loaded %d cached facts for package %q and analyzer %s", len(facts), act.Package.Name, act.Analyzer.Name)

	var importsByPath map[string]*types.Package

	// Lazily built lookup of the package owning each fact (by import path),
	// resolved through this package's transitive imports.
	resolvePkg := func(pkgPath string) *types.Package {
		if pkgPath == "" {
			return act.Package.Types
		}

		if importsByPath == nil {
			importsByPath = collectImports(act.Package.Types)
		}

		return importsByPath[pkgPath]
	}

	for _, f := range facts {
		pkg := resolvePkg(f.PkgPath)
		if pkg == nil {
			// The owning package is not reachable from this package anymore.
			continue
		}

		if f.Path == "" { // this is a package fact
			key := packageFactKey{pkg: pkg, typ: act.factType(f.Fact)}
			act.packageFacts[key] = f.Fact
			continue
		}

		obj, err := objectpath.Object(pkg, objectpath.Path(f.Path))
		if err != nil {
			// Be lenient about these errors.
			// For example, when analyzing io/ioutil from source,
			// we may get a fact for methods on the devNull type,
			// and objectpath will happily create a path for them.
			// However,
			// when we later load io/ioutil from export data,
			// the path no longer resolves.
			//
			// If an exported type embeds the unexported type,
			// then (part of) the unexported type will become part of the type information and our path will resolve again.
			continue
		}

		factKey := objectFactKey{obj, act.factType(f.Fact)}

		act.objectFacts[factKey] = f.Fact
	}

	return true
}

// collectImports returns a map of every package transitively imported by pkg (including pkg itself), keyed by import path.
// It is used to resolve facts about objects that belong to a dependency reachable through the export data.
func collectImports(pkg *types.Package) map[string]*types.Package {
	result := map[string]*types.Package{pkg.Path(): pkg}

	var visit func(p *types.Package)

	visit = func(p *types.Package) {
		for _, imp := range p.Imports() {
			if _, ok := result[imp.Path()]; ok {
				continue
			}

			result[imp.Path()] = imp
			visit(imp)
		}
	}

	visit(pkg)

	return result
}

func factCacheKey(a *analysis.Analyzer) string {
	return fmt.Sprintf("%s/facts", a.Name)
}
