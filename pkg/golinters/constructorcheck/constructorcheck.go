// Package analyzer is a linter that reports ignored constructors.
// It shows you places where someone is doing T{} or &T{}
// instead of using NewT declared in the same package as T.
// A constructor for type T (only structs are supported at the moment)
// is a function with name "NewT" that returns a value of type T or *T.
// Types returned by constructors are not checked right now,
// only that type T inferred from the function name exists in the same package.
// Standard library packages are excluded from analysis.
package constructorcheck

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type ConstructorFact struct {
	ConstructorName string
	Pos             token.Pos
	End             token.Pos
}

func (f *ConstructorFact) AFact() {}

var Analyzer = &analysis.Analyzer{
	Name:      "constructor_check",
	Doc:       "check for types constructed manually ignoring constructor",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{(*ConstructorFact)(nil)},
}

var stdPackages = stdPackageNames()

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.FuncDecl)(nil),
	}

	zeroValues := make(map[token.Pos]types.Object)
	nilValues := make(map[token.Pos]types.Object)
	compositeLiterals := make(map[token.Pos]types.Object)
	typeAliases := make(map[types.Object]types.Object)

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		switch decl := node.(type) {
		case *ast.CallExpr:
			// check if it's a new call
			fn, ok := decl.Fun.(*ast.Ident)
			if !ok {
				break
			}
			if fn.Name != "new" {
				break
			}
			// check we have only one argument (the type)
			if len(decl.Args) != 1 {
				break
			}

			ident := typeIdent(decl.Args[0])
			if ident == nil {
				break
			}

			typeObj := pass.TypesInfo.ObjectOf(ident)
			if typeObj == nil {
				break
			}
			zeroValues[node.Pos()] = typeObj
		case *ast.ValueSpec:
			// check it's a pointer value
			starExpr, ok := decl.Type.(*ast.StarExpr)
			if !ok {
				break
			}
			// check it's using a named type
			ident := typeIdent(starExpr.X)
			if ident == nil {
				break
			}
			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}
			nilValues[node.Pos()] = obj
		case *ast.CompositeLit:
			ident := typeIdent(decl.Type)
			if ident == nil {
				break
			}

			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}
			// if it's a zero value literal
			if decl.Elts == nil {
				zeroValues[node.Pos()] = obj
				break
			}
			compositeLiterals[node.Pos()] = obj
		case *ast.TypeSpec:
			// get base type if any
			baseIdent := typeIdent(decl.Type)
			if baseIdent == nil {
				break
			}
			// get base type object
			baseTypeObj := pass.TypesInfo.ObjectOf(baseIdent)
			if baseTypeObj == nil {
				break
			}

			// get this type's object
			typeObj := pass.TypesInfo.ObjectOf(decl.Name)
			if typeObj == nil {
				break
			}
			typeAliases[typeObj] = baseTypeObj
		case *ast.FuncDecl:
			// check if it's a function not a method
			if decl.Recv != nil {
				break
			}

			// check if function name starts with "New"
			if !strings.HasPrefix(decl.Name.Name, "New") {
				break
			}

			// check if function name follows the NewT template
			// TODO: think about easing this requirement because often
			// they rename types and forget to rename constructors
			typeName, ok := strings.CutPrefix(decl.Name.Name, "New")
			if !ok {
				break
			}

			// check if type T extracted from function name exists
			obj := pass.Pkg.Scope().Lookup(typeName)
			if obj == nil {
				break
			}

			// ignore standard library types
			if _, ok := stdPackages[obj.Pkg().Name()]; ok {
				break
			}
			// check if supposed constructor returns exactly one value
			// TODO: implement other cases ?
			// (T, err), (*T, err), (T, bool), (*T, bool)
			returns := decl.Type.Results.List
			if len(returns) != 1 {
				break
			}
			// to be done later:
			// // check if supposed constructor returns a value of type T or *T
			// // declared in the same package and T equals extracted type name

			// assume we have a valid constructor
			fact := ConstructorFact{
				ConstructorName: decl.Name.Name,
				Pos:             decl.Pos(),
				End:             decl.End(),
			}
			pass.ExportObjectFact(obj, &fact)
		default:
			// fmt.Printf("%#v\n", node)
		}
	})

	for typeObj, baseTypeObj := range typeAliases {
		// check the base type has a constructor
		existingFact := new(ConstructorFact)
		if !pass.ImportObjectFact(baseTypeObj, existingFact) {
			continue
		}

		// mark derived type as having constructor
		newFact := ConstructorFact{
			ConstructorName: existingFact.ConstructorName,
			Pos:             existingFact.Pos,
			End:             existingFact.End,
		}
		pass.ExportObjectFact(typeObj, &newFact)
	}

	for pos, obj := range nilValues {
		if constr, ok := constructorName(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"nil value of type %s may be unsafe, use constructor %s instead",
				obj.Type(),
				constr,
			)
		}
	}
	for pos, obj := range zeroValues {
		if constr, ok := constructorName(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"zero value of type %s may be unsafe, use constructor %s instead",
				obj.Type(),
				constr,
			)
		}
	}
	for pos, obj := range compositeLiterals {
		if constr, ok := constructorName(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"use constructor %s for type %s instead of a composite literal",
				constr,
				obj.Type(),
			)
		}
	}

	return nil, nil
}

func constructorName(pass *analysis.Pass, obj types.Object, pos token.Pos) (string, bool) {
	fact := new(ConstructorFact)
	if !pass.ImportObjectFact(obj, fact) {
		return "", false
	}

	// if used inside T's constructor - ignore
	if pos >= fact.Pos &&
		pos < fact.End {
		return "", false
	}

	return fact.ConstructorName, true
}

// typeIdent returns either local or imported type ident or nil
func typeIdent(expr ast.Expr) *ast.Ident {
	switch id := expr.(type) {
	case *ast.Ident:
		return id
	case *ast.SelectorExpr:
		return id.Sel
	}
	return nil
}

func stdPackageNames() map[string]struct{} {
	// inspired by https://pkg.go.dev/golang.org/x/tools/go/packages#Load
	cmd := exec.Command("go", "list", "std")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal("can't load standard library package names")
	}
	pkgs := strings.Fields(string(output))

	stdPkgNames := make(map[string]struct{}, len(pkgs))
	for _, pkg := range pkgs {
		stdPkgNames[pkg] = struct{}{}
	}
	return stdPkgNames
}
