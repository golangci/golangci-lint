package packages

import (
	"go/build"
	"path/filepath"
)

type Package struct {
	bp *build.Package

	isFake bool
}

func (pkg *Package) Files(includeTest bool) []string {
	pkgFiles := append([]string{}, pkg.bp.GoFiles...)

	// TODO: add cgo files
	if includeTest {
		pkgFiles = append(pkgFiles, pkg.bp.TestGoFiles...)
		pkgFiles = append(pkgFiles, pkg.bp.XTestGoFiles...)
	}

	for i, f := range pkgFiles {
		pkgFiles[i] = filepath.Join(pkg.bp.Dir, f)
	}

	return pkgFiles
}
