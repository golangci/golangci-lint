package astcache

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

type File struct {
	F    *ast.File
	Fset *token.FileSet
	Name string
	Err  error
}

type Cache struct {
	m   map[string]*File // map from absolute file path to file data
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

func (c Cache) keys() []string {
	var keys []string
	for k := range c.m {
		keys = append(keys, k)
	}
	return keys
}

func (c Cache) GetOrParse(filename string, fset *token.FileSet) *File {
	if !filepath.IsAbs(filename) {
		absFilename, err := filepath.Abs(filename)
		if err != nil {
			c.log.Warnf("Can't abs-ify filename %s: %s", filename, err)
		} else {
			filename = absFilename
		}
	}

	f := c.m[filename]
	if f != nil {
		return f
	}

	c.log.Infof("Parse AST for file %s on demand, existing files are %s",
		filename, strings.Join(c.keys(), ","))
	c.parseFile(filename, fset)
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

func LoadFromPackages(pkgs []*packages.Package, log logutils.Log) (*Cache, error) {
	c := NewCache(log)

	for _, pkg := range pkgs {
		c.loadFromPackage(pkg)
	}

	c.prepareValidFiles()
	return c, nil
}

func (c *Cache) loadFromPackage(pkg *packages.Package) {
	if len(pkg.Syntax) == 0 || len(pkg.GoFiles) != len(pkg.CompiledGoFiles) {
		// len(pkg.Syntax) == 0 if only filenames are loaded
		// lengths aren't equal if there are preprocessed files (cgo)
		startedAt := time.Now()

		// can't use pkg.Fset: it will overwrite offsets by preprocessed files
		fset := token.NewFileSet()
		for _, f := range pkg.GoFiles {
			c.parseFile(f, fset)
		}

		c.log.Infof("Parsed AST of all pkg.GoFiles: %s for %s", pkg.GoFiles, time.Since(startedAt))
		return
	}

	for _, f := range pkg.Syntax {
		pos := pkg.Fset.Position(f.Pos())
		if pos.Filename == "" {
			continue
		}

		c.m[pos.Filename] = &File{
			F:    f,
			Fset: pkg.Fset,
			Name: pos.Filename,
		}
	}
}

func (c *Cache) parseFile(filePath string, fset *token.FileSet) {
	if fset == nil {
		fset = token.NewFileSet()
	}

	// comments needed by e.g. golint
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
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
