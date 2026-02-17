package goanalysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func Test_action_markDepsForAnalyzingSource(t *testing.T) {
	t.Run("marks direct horizontal deps", func(t *testing.T) {
		pkg := &packages.Package{PkgPath: "pkg/a"}

		dep := &action{Package: pkg}
		act := &action{
			Package: pkg,
			Deps:    []*action{dep},
		}

		act.markDepsForAnalyzingSource()

		assert.True(t, dep.needAnalyzeSource)
	})

	t.Run("marks transitive horizontal deps", func(t *testing.T) {
		pkg := &packages.Package{PkgPath: "pkg/a"}

		// Chain: act -> dep1 -> dep2
		dep2 := &action{Package: pkg}
		dep1 := &action{
			Package: pkg,
			Deps:    []*action{dep2},
		}
		act := &action{
			Package: pkg,
			Deps:    []*action{dep1},
		}

		act.markDepsForAnalyzingSource()

		assert.True(t, dep1.needAnalyzeSource)
		assert.True(t, dep2.needAnalyzeSource)
	})

	t.Run("marks deep transitive horizontal deps", func(t *testing.T) {
		pkg := &packages.Package{PkgPath: "pkg/a"}

		// Chain: act -> dep1 -> dep2 -> dep3 (simulates nilaway -> buildssa -> ctrlflow -> inspect)
		dep3 := &action{Package: pkg}
		dep2 := &action{
			Package: pkg,
			Deps:    []*action{dep3},
		}
		dep1 := &action{
			Package: pkg,
			Deps:    []*action{dep2},
		}
		act := &action{
			Package: pkg,
			Deps:    []*action{dep1},
		}

		act.markDepsForAnalyzingSource()

		assert.True(t, dep1.needAnalyzeSource)
		assert.True(t, dep2.needAnalyzeSource)
		assert.True(t, dep3.needAnalyzeSource)
	})

	t.Run("does not mark cross-package deps", func(t *testing.T) {
		pkgA := &packages.Package{PkgPath: "pkg/a"}
		pkgB := &packages.Package{PkgPath: "pkg/b"}

		dep := &action{Package: pkgB}
		act := &action{
			Package: pkgA,
			Deps:    []*action{dep},
		}

		act.markDepsForAnalyzingSource()

		assert.False(t, dep.needAnalyzeSource)
	})

	t.Run("handles cycles without infinite recursion", func(t *testing.T) {
		pkg := &packages.Package{PkgPath: "pkg/a"}

		dep1 := &action{Package: pkg}
		dep2 := &action{Package: pkg}

		dep1.Deps = []*action{dep2}
		dep2.Deps = []*action{dep1}

		act := &action{
			Package: pkg,
			Deps:    []*action{dep1},
		}

		// Should not hang or panic
		act.markDepsForAnalyzingSource()

		assert.True(t, dep1.needAnalyzeSource)
		assert.True(t, dep2.needAnalyzeSource)
	})

	t.Run("skips already marked deps", func(t *testing.T) {
		pkg := &packages.Package{PkgPath: "pkg/a"}

		dep := &action{
			Package:           pkg,
			needAnalyzeSource: true, // already marked
		}
		act := &action{
			Package: pkg,
			Deps:    []*action{dep},
		}

		// Should not recurse into already-marked dep
		act.markDepsForAnalyzingSource()

		assert.True(t, dep.needAnalyzeSource)
	})
}
