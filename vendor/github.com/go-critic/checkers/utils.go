package checkers

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/go-critic/checkers/internal/lintutil"
	"github.com/go-lintpack/lintpack"
	"golang.org/x/tools/go/ast/astutil"
)

// isUnitTestFunc reports whether FuncDecl declares testing function.
func isUnitTestFunc(ctx *lintpack.CheckerContext, fn *ast.FuncDecl) bool {
	if !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	typ := ctx.TypesInfo.TypeOf(fn.Name)
	if sig, ok := typ.(*types.Signature); ok {
		return sig.Results().Len() == 0 &&
			sig.Params().Len() == 1 &&
			sig.Params().At(0).Type().String() == "*testing.T"
	}
	return false
}

// qualifiedName returns called expr fully-quallified name.
//
// It works for simple identifiers like f => "f" and identifiers
// from other package like pkg.f => "pkg.f".
//
// For all unexpected expressions returns empty string.
func qualifiedName(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.SelectorExpr:
		pkg, ok := x.X.(*ast.Ident)
		if !ok {
			return ""
		}
		return pkg.Name + "." + x.Sel.Name
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}

// identOf returns identifier for x that can be used to obtain associated types.Object.
// Returns nil for expressions that yield temporary results, like `f().field`.
func identOf(x ast.Node) *ast.Ident {
	switch x := x.(type) {
	case *ast.Ident:
		return x
	case *ast.SelectorExpr:
		return identOf(x.Sel)
	case *ast.TypeAssertExpr:
		// x.(type) - x may contain ident.
		return identOf(x.X)
	case *ast.IndexExpr:
		// x[i] - x may contain ident.
		return identOf(x.X)
	case *ast.StarExpr:
		// *x - x may contain ident.
		return identOf(x.X)
	case *ast.SliceExpr:
		// x[:] - x may contain ident.
		return identOf(x.X)

	default:
		// Note that this function is not comprehensive.
		return nil
	}
}

// findNode applies pred for root and all it's childs until it returns true.
// Matched node is returned.
// If none of the nodes matched predicate, nil is returned.
func findNode(root ast.Node, pred func(ast.Node) bool) ast.Node {
	var found ast.Node
	astutil.Apply(root, nil, func(cur *astutil.Cursor) bool {
		if pred(cur.Node()) {
			found = cur.Node()
			return false
		}
		return true
	})
	return found
}

// containsNode reports whether `findNode(root, pred)!=nil`.
func containsNode(root ast.Node, pred func(ast.Node) bool) bool {
	return findNode(root, pred) != nil
}

// isSafeExpr reports whether expr is softly safe expression and contains
// no significant side-effects. As opposed to strictly safe expressions,
// soft safe expressions permit some forms of side-effects, like
// panic possibility during indexing or nil pointer dereference.
//
// Uses types info to determine type conversion expressions that
// are the only permitted kinds of call expressions.
func isSafeExpr(info *types.Info, expr ast.Expr) bool {
	// This list switch is not comprehensive and uses
	// whitelist to be on the conservative side.
	// Can be extended as needed.
	//
	// Note that it is not very strict "safe" as
	// index expressions are permitted even though they
	// may cause panics.
	switch expr := expr.(type) {
	case *ast.StarExpr:
		return isSafeExpr(info, expr.X)
	case *ast.BinaryExpr:
		return isSafeExpr(info, expr.X) && isSafeExpr(info, expr.Y)
	case *ast.UnaryExpr:
		return expr.Op != token.ARROW && isSafeExpr(info, expr.X)
	case *ast.BasicLit, *ast.Ident:
		return true
	case *ast.IndexExpr:
		return isSafeExpr(info, expr.X) && isSafeExpr(info, expr.Index)
	case *ast.SelectorExpr:
		return isSafeExpr(info, expr.X)
	case *ast.ParenExpr:
		return isSafeExpr(info, expr.X)
	case *ast.CompositeLit:
		return isSafeExprList(info, expr.Elts)
	case *ast.CallExpr:
		return lintutil.IsTypeExpr(info, expr.Fun) &&
			isSafeExprList(info, expr.Args)

	default:
		return false
	}
}

// isSafeExprList reports whether every expr in list is safe.
// See isSafeExpr.
func isSafeExprList(info *types.Info, list []ast.Expr) bool {
	for _, expr := range list {
		if !isSafeExpr(info, expr) {
			return false
		}
	}
	return true
}
