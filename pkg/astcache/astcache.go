package astcache

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/loader"
)

type File struct {
	F    *ast.File
	Fset *token.FileSet
	err  error
}

type Cache struct {
	m map[string]*File
	s []*File
}

func (c Cache) GetAllValidFiles() []*File {
	return c.s
}

func (c *Cache) prepareValidFiles() {
	files := make([]*File, 0, len(c.m))
	for _, f := range c.m {
		if f.err != nil || f.F == nil {
			continue
		}
		files = append(files, f)
	}
	c.s = files
}

func LoadFromProgram(prog *loader.Program) *Cache {
	c := &Cache{
		m: map[string]*File{},
	}
	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			pos := prog.Fset.Position(0)
			c.m[pos.Filename] = &File{
				F:    f,
				Fset: prog.Fset,
			}
		}
	}

	c.prepareValidFiles()
	return c
}

func LoadFromFiles(files []string) *Cache {
	c := &Cache{
		m: map[string]*File{},
	}
	fset := token.NewFileSet()
	for _, filePath := range files {
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments) // comments needed by e.g. golint
		c.m[filePath] = &File{
			F:    f,
			Fset: fset,
			err:  err,
		}
	}

	c.prepareValidFiles()
	return c
}
