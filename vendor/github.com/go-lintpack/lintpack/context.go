package lintpack

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/go-toolsmith/astfmt"
)

// Context is a readonly state shared among every checker.
type Context struct {
	// TypesInfo carries parsed packages types information.
	TypesInfo *types.Info

	// SizesInfo carries alignment and type size information.
	// Arch-dependent.
	SizesInfo types.Sizes

	// FileSet is a file set that was used during the program loading.
	FileSet *token.FileSet

	// Pkg describes package that is being checked.
	Pkg *types.Package

	// Filename is a currently checked file name.
	Filename string

	// Require records what optional resources are required
	// by the checkers set that use this context.
	//
	// Every require fields makes associated context field
	// to be properly initialized.
	// For example, Context.require.PkgObjects => Context.PkgObjects.
	Require struct {
		PkgObjects bool
		PkgRenames bool
	}

	// PkgObjects stores all imported packages and their local names.
	PkgObjects map[*types.PkgName]string

	// PkgRenames maps package path to its local renaming.
	// Contains no entries for packages that were imported without
	// explicit local names.
	PkgRenames map[string]string
}

// NewContext returns new shared context to be used by every checker.
//
// All data carried by the context is readonly for checkers,
// but can be modified by the integrating application.
func NewContext(fset *token.FileSet, sizes types.Sizes) *Context {
	return &Context{
		FileSet:   fset,
		SizesInfo: sizes,
		TypesInfo: &types.Info{},
	}
}

// SetPackageInfo sets package-related metadata.
//
// Must be called for every package being checked.
func (c *Context) SetPackageInfo(info *types.Info, pkg *types.Package) {
	if info != nil {
		// We do this kind of assignment to avoid
		// changing c.typesInfo field address after
		// every re-assignment.
		*c.TypesInfo = *info
	}
	c.Pkg = pkg
}

// SetFileInfo sets file-related metadata.
//
// Must be called for every source code file being checked.
func (c *Context) SetFileInfo(name string, f *ast.File) {
	c.Filename = name
	if c.Require.PkgObjects {
		resolvePkgObjects(c, f)
	}
	if c.Require.PkgRenames {
		resolvePkgRenames(c, f)
	}
}

// CheckerContext is checker-local context copy.
// Fields that are not from Context itself are writeable.
type CheckerContext struct {
	*Context

	// printer used to format warning text.
	printer *astfmt.Printer

	Params parameters

	warnings []Warning
}

// Warn adds a Warning to checker output.
func (ctx *CheckerContext) Warn(node ast.Node, format string, args ...interface{}) {
	ctx.warnings = append(ctx.warnings, Warning{
		Text: ctx.printer.Sprintf(format, args...),
		Node: node,
	})
}

func resolvePkgObjects(ctx *Context, f *ast.File) {
	ctx.PkgObjects = make(map[*types.PkgName]string, len(f.Imports))

	for _, spec := range f.Imports {
		if spec.Name != nil {
			obj := ctx.TypesInfo.ObjectOf(spec.Name)
			ctx.PkgObjects[obj.(*types.PkgName)] = spec.Name.Name
		} else {
			obj := ctx.TypesInfo.Implicits[spec]
			ctx.PkgObjects[obj.(*types.PkgName)] = obj.Name()
		}
	}
}

func resolvePkgRenames(ctx *Context, f *ast.File) {
	ctx.PkgRenames = make(map[string]string)

	for _, spec := range f.Imports {
		if spec.Name != nil {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				panic(err)
			}
			ctx.PkgRenames[path] = spec.Name.Name
		}
	}
}
