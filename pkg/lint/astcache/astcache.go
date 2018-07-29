package astcache

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/goutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"golang.org/x/tools/go/loader"
)

type File struct {
	F    *ast.File
	Fset *token.FileSet
	Name string
	Err  error
}

type Cache struct {
	m   map[string]*File
	s   []*File
	log logutils.Log
}

func NewCache(log logutils.Log) *Cache {
	return &Cache{
		m:   map[string]*File{},
		log: log,
	}
}

func (c Cache) Get(filename string) *File {
	return c.m[filepath.Clean(filename)]
}

func (c Cache) GetOrParse(filename string) *File {
	f := c.m[filename]
	if f != nil {
		return f
	}

	c.log.Infof("Parse AST for file %s on demand", filename)
	c.parseFile(filename, nil)
	return c.m[filename]
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

func LoadFromProgram(prog *loader.Program, log logutils.Log) (*Cache, error) {
	c := NewCache(log)

	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			pos := prog.Fset.Position(f.Pos())
			if pos.Filename == "" {
				continue
			}

			if goutils.IsCgoFilename(pos.Filename) {
				continue
			}

			path, err := fsutils.ShortestRelPath(pos.Filename, "")
			if err != nil {
				c.log.Warnf("Can't get relative path for %s: %s",
					pos.Filename, err)
				continue
			}

			c.m[path] = &File{
				F:    f,
				Fset: prog.Fset,
				Name: path,
			}
		}
	}

	c.prepareValidFiles()
	return c, nil
}

func (c *Cache) parseFile(filePath string, fset *token.FileSet) {
	if fset == nil {
		fset = token.NewFileSet()
	}

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments) // comments needed by e.g. golint
	c.m[filePath] = &File{
		F:    f,
		Fset: fset,
		Err:  err,
		Name: filePath,
	}
	if err != nil {
		c.log.Warnf("Can't parse AST of %s: %s", filePath, err)
	}
}

func LoadFromFiles(files []string, log logutils.Log) (*Cache, error) { //nolint:unparam
	c := NewCache(log)

	fset := token.NewFileSet()
	for _, filePath := range files {
		filePath = filepath.Clean(filePath)
		c.parseFile(filePath, fset)
	}

	c.prepareValidFiles()
	return c, nil
}
