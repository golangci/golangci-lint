// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gas

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

// MatchCallByPackage ensures that the specified package is imported,
// adjusts the name for any aliases and ignores cases that are
// initialization only imports.
//
// Usage:
// 	node, matched := MatchCallByPackage(n, ctx, "math/rand", "Read")
//
func MatchCallByPackage(n ast.Node, c *Context, pkg string, names ...string) (*ast.CallExpr, bool) {

	importedName, found := GetImportedName(pkg, c)
	if !found {
		return nil, false
	}

	if callExpr, ok := n.(*ast.CallExpr); ok {
		packageName, callName, err := GetCallInfo(callExpr, c)
		if err != nil {
			return nil, false
		}
		if packageName == importedName {
			for _, name := range names {
				if callName == name {
					return callExpr, true
				}
			}
		}
	}
	return nil, false
}

// MatchCallByType ensures that the node is a call expression to a
// specific object type.
//
// Usage:
// 	node, matched := MatchCallByType(n, ctx, "bytes.Buffer", "WriteTo", "Write")
//
func MatchCallByType(n ast.Node, ctx *Context, requiredType string, calls ...string) (*ast.CallExpr, bool) {
	if callExpr, ok := n.(*ast.CallExpr); ok {
		typeName, callName, err := GetCallInfo(callExpr, ctx)
		if err != nil {
			return nil, false
		}
		if typeName == requiredType {
			for _, call := range calls {
				if call == callName {
					return callExpr, true
				}
			}
		}
	}
	return nil, false
}

// MatchCompLit will match an ast.CompositeLit based on the supplied type
func MatchCompLit(n ast.Node, ctx *Context, required string) *ast.CompositeLit {
	if complit, ok := n.(*ast.CompositeLit); ok {
		typeOf := ctx.Info.TypeOf(complit)
		if typeOf.String() == required {
			return complit
		}
	}
	return nil
}

// GetInt will read and return an integer value from an ast.BasicLit
func GetInt(n ast.Node) (int64, error) {
	if node, ok := n.(*ast.BasicLit); ok && node.Kind == token.INT {
		return strconv.ParseInt(node.Value, 0, 64)
	}
	return 0, fmt.Errorf("Unexpected AST node type: %T", n)
}

// GetFloat will read and return a float value from an ast.BasicLit
func GetFloat(n ast.Node) (float64, error) {
	if node, ok := n.(*ast.BasicLit); ok && node.Kind == token.FLOAT {
		return strconv.ParseFloat(node.Value, 64)
	}
	return 0.0, fmt.Errorf("Unexpected AST node type: %T", n)
}

// GetChar will read and return a char value from an ast.BasicLit
func GetChar(n ast.Node) (byte, error) {
	if node, ok := n.(*ast.BasicLit); ok && node.Kind == token.CHAR {
		return node.Value[0], nil
	}
	return 0, fmt.Errorf("Unexpected AST node type: %T", n)
}

// GetString will read and return a string value from an ast.BasicLit
func GetString(n ast.Node) (string, error) {
	if node, ok := n.(*ast.BasicLit); ok && node.Kind == token.STRING {
		return strconv.Unquote(node.Value)
	}
	return "", fmt.Errorf("Unexpected AST node type: %T", n)
}

// GetCallObject returns the object and call expression and associated
// object for a given AST node. nil, nil will be returned if the
// object cannot be resolved.
func GetCallObject(n ast.Node, ctx *Context) (*ast.CallExpr, types.Object) {
	switch node := n.(type) {
	case *ast.CallExpr:
		switch fn := node.Fun.(type) {
		case *ast.Ident:
			return node, ctx.Info.Uses[fn]
		case *ast.SelectorExpr:
			return node, ctx.Info.Uses[fn.Sel]
		}
	}
	return nil, nil
}

// GetCallInfo returns the package or type and name  associated with a
// call expression.
func GetCallInfo(n ast.Node, ctx *Context) (string, string, error) {
	switch node := n.(type) {
	case *ast.CallExpr:
		switch fn := node.Fun.(type) {
		case *ast.SelectorExpr:
			switch expr := fn.X.(type) {
			case *ast.Ident:
				if expr.Obj != nil && expr.Obj.Kind == ast.Var {
					t := ctx.Info.TypeOf(expr)
					if t != nil {
						return t.String(), fn.Sel.Name, nil
					}
					return "undefined", fn.Sel.Name, fmt.Errorf("missing type info")
				}
				return expr.Name, fn.Sel.Name, nil
			}
		case *ast.Ident:
			return ctx.Pkg.Name(), fn.Name, nil
		}
	}
	return "", "", fmt.Errorf("unable to determine call info")
}

// GetImportedName returns the name used for the package within the
// code. It will resolve aliases and ignores initalization only imports.
func GetImportedName(path string, ctx *Context) (string, bool) {
	importName, imported := ctx.Imports.Imported[path]
	if !imported {
		return "", false
	}

	if _, initonly := ctx.Imports.InitOnly[path]; initonly {
		return "", false
	}

	if alias, ok := ctx.Imports.Aliased[path]; ok {
		importName = alias
	}
	return importName, true
}

// GetImportPath resolves the full import path of an identifer based on
// the imports in the current context.
func GetImportPath(name string, ctx *Context) (string, bool) {
	for path := range ctx.Imports.Imported {
		if imported, ok := GetImportedName(path, ctx); ok && imported == name {
			return path, true
		}
	}
	return "", false
}

// GetLocation returns the filename and line number of an ast.Node
func GetLocation(n ast.Node, ctx *Context) (string, int) {
	fobj := ctx.FileSet.File(n.Pos())
	return fobj.Name(), fobj.Line(n.Pos())
}
