package testutils

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/securego/gosec/v2"
)

type buildObj struct {
	pkg    *build.Package
	config *packages.Config
	pkgs   []*packages.Package
}

// TestPackage is a mock package for testing purposes
type TestPackage struct {
	Path   string
	Files  map[string]string
	onDisk bool
	build  *buildObj
}

// Option provides a way to adjust the package config depending on testing
// requirements
type Option func(conf *packages.Config)

// WithBuildTags enables injecting build tags into the package config on build.
func WithBuildTags(tags []string) Option {
	return func(conf *packages.Config) {
		conf.BuildFlags = tags
	}
}

// NewTestPackage will create a new and empty package. Must call Close() to cleanup
// auxiliary files
func NewTestPackage() *TestPackage {
	workingDir, err := os.MkdirTemp("", "gosecs_test")
	if err != nil {
		return nil
	}

	return &TestPackage{
		Path:   workingDir,
		Files:  make(map[string]string),
		onDisk: false,
		build:  nil,
	}
}

// AddFile inserts the filename and contents into the package contents
func (p *TestPackage) AddFile(filename, content string) {
	p.Files[path.Join(p.Path, filename)] = content
}

func (p *TestPackage) write() error {
	if p.onDisk {
		return nil
	}
	for filename, content := range p.Files {
		if e := os.WriteFile(filename, []byte(content), 0o644); e != nil /* #nosec G306 */ {
			return e
		}
	}
	p.onDisk = true
	return nil
}

// Build ensures all files are persisted to disk and built
func (p *TestPackage) Build(opts ...Option) error {
	if p.build != nil {
		return nil
	}
	if err := p.write(); err != nil {
		return err
	}

	conf := &packages.Config{
		Mode:  gosec.LoadMode,
		Tests: false,
	}
	for _, opt := range opts {
		opt(conf)
	}

	// step 1/2: build context requires the array of build tags.
	builder := build.Default
	builder.BuildTags = conf.BuildFlags
	basePackage, err := builder.ImportDir(p.Path, build.ImportComment)
	if err != nil {
		return err
	}

	var packageFiles []string
	for _, filename := range basePackage.GoFiles {
		packageFiles = append(packageFiles, path.Join(p.Path, filename))
	}

	// step 2/2: normalise to cli build flags for package loading
	conf.BuildFlags = gosec.CLIBuildTags(conf.BuildFlags)
	pkgs, err := packages.Load(conf, packageFiles...)
	if err != nil {
		return err
	}
	p.build = &buildObj{
		pkg:    basePackage,
		config: conf,
		pkgs:   pkgs,
	}
	return nil
}

// CreateContext builds a context out of supplied package context
func (p *TestPackage) CreateContext(filename string, opts ...Option) *gosec.Context {
	if err := p.Build(opts...); err != nil {
		log.Fatal(err)
		return nil
	}

	for _, pkg := range p.build.pkgs {
		for _, file := range pkg.Syntax {
			pkgFile := pkg.Fset.File(file.Pos()).Name()
			strip := fmt.Sprintf("%s%c", p.Path, os.PathSeparator)
			pkgFile = strings.TrimPrefix(pkgFile, strip)
			if pkgFile == filename {
				ctx := &gosec.Context{
					FileSet:      pkg.Fset,
					Root:         file,
					Config:       gosec.NewConfig(),
					Info:         pkg.TypesInfo,
					Pkg:          pkg.Types,
					Imports:      gosec.NewImportTracker(),
					PassedValues: make(map[string]interface{}),
				}
				ctx.Imports.TrackPackages(ctx.Pkg.Imports()...)
				return ctx
			}
		}
	}
	return nil
}

// Close will delete the package and all files in that directory
func (p *TestPackage) Close() {
	if p.onDisk {
		err := os.RemoveAll(p.Path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Pkgs returns the current built packages
func (p *TestPackage) Pkgs() []*packages.Package {
	if p.build != nil {
		return p.build.pkgs
	}
	return []*packages.Package{}
}

// PrintErrors prints to os.Stderr the accumulated errors of built packages
func (p *TestPackage) PrintErrors() int {
	return packages.PrintErrors(p.Pkgs())
}
