package packages

import (
	"fmt"
	"go/build"
)

type Program struct {
	packages []Package

	bctx build.Context
}

func (p *Program) String() string {
	files := p.Files(true)
	if len(files) == 1 {
		return files[0]
	}

	return fmt.Sprintf("%s", p.Dirs())
}

func (p *Program) BuildContext() *build.Context {
	return &p.bctx
}

func (p Program) Packages() []Package {
	return p.packages
}

func (p *Program) addPackage(pkg *Package) {
	packagesToAdd := []Package{*pkg}
	if len(pkg.bp.XTestGoFiles) != 0 {
		// create separate package because xtest files have different package name
		xbp := build.Package{
			Dir:            pkg.bp.Dir,
			ImportPath:     pkg.bp.ImportPath + "_test",
			XTestGoFiles:   pkg.bp.XTestGoFiles,
			XTestImportPos: pkg.bp.XTestImportPos,
			XTestImports:   pkg.bp.XTestImports,
		}
		packagesToAdd = append(packagesToAdd, Package{
			bp: &xbp,
		})
		pkg.bp.XTestGoFiles = nil
		pkg.bp.XTestImportPos = nil
		pkg.bp.XTestImports = nil
	}

	p.packages = append(p.packages, packagesToAdd...)
}

func (p *Program) Files(includeTest bool) []string {
	var ret []string
	for _, pkg := range p.packages {
		ret = append(ret, pkg.Files(includeTest)...)
	}

	return ret
}

func (p *Program) Dirs() []string {
	var ret []string
	for _, pkg := range p.packages {
		if !pkg.isFake {
			ret = append(ret, pkg.Dir())
		}
	}

	return ret
}
