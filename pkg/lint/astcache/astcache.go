package astcache

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
)

type File struct {
	F    *ast.File
	Fset *token.FileSet
	Name string
	Err  error
}

type Cache struct {
	m map[string]*File
	s []*File
}

func (c Cache) Get(filename string) *File {
	return c.m[filepath.Clean(filename)]
}

func (c Cache) GetAllValidFiles() []*File {
	return c.s
}

func (c *Cache) prepareValidFiles() {
	files := make([]*File, 0, len(c.m))
	for _, f := range c.m {
		if f.Err != nil || f.F == nil {
			continue
		}
		files = append(files, f)
	}
	c.s = files
}

func LoadFromProgram(prog *loader.Program) (*Cache, error) {
	c := &Cache{
		m: map[string]*File{},
	}

	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't get working dir: %s", err)
	}

	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			pos := prog.Fset.Position(f.Pos())
			if pos.Filename == "" {
				continue
			}

			relPath, err := filepath.Rel(root, pos.Filename)
			if err != nil {
				logrus.Warnf("Can't get relative path for %s and %s: %s",
					root, pos.Filename, err)
				continue
			}

			c.m[relPath] = &File{
				F:    f,
				Fset: prog.Fset,
				Name: relPath,
			}
		}
	}

	c.prepareValidFiles()
	return c, nil
}

func LoadFromFiles(files []string) (*Cache, error) {
	c := &Cache{
		m: map[string]*File{},
	}

	fset := token.NewFileSet()
	for _, filePath := range files {
		filePath = filepath.Clean(filePath)

		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments) // comments needed by e.g. golint
		c.m[filePath] = &File{
			F:    f,
			Fset: fset,
			Err:  err,
			Name: filePath,
		}
	}

	c.prepareValidFiles()
	return c, nil
}
