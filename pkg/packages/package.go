package packages

import (
	"go/build"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/goutils"
)

type Package struct {
	bp *build.Package

	isFake bool
	dir    string // dir != bp.dir only if isFake == true
}

func (pkg *Package) Files(includeTest bool) []string {
	var pkgFiles []string
	for _, f := range pkg.bp.GoFiles {
		if !goutils.IsCgoFilename(f) {
			// skip cgo at all levels to prevent failures on file reading
			pkgFiles = append(pkgFiles, f)
		}
	}

	// TODO: add cgo files
	if includeTest {
		pkgFiles = append(pkgFiles, pkg.TestFiles()...)
	}

	for i, f := range pkgFiles {
		pkgFiles[i] = filepath.Join(pkg.bp.Dir, f)
	}

	return pkgFiles
}

func (pkg *Package) Dir() string {
	if pkg.dir != "" { // for fake packages
		return pkg.dir
	}

	return pkg.bp.Dir
}

func (pkg *Package) IsTestOnly() bool {
	return len(pkg.bp.GoFiles) == 0
}

func (pkg *Package) TestFiles() []string {
	var pkgFiles []string
	pkgFiles = append(pkgFiles, pkg.bp.TestGoFiles...)
	pkgFiles = append(pkgFiles, pkg.bp.XTestGoFiles...)
	return pkgFiles
}
