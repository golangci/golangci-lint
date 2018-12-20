package checkers

import (
	"go/ast"
	"go/types"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/astcopy"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "methodExprCall"
	info.Tags = []string{"style", "experimental"}
	info.Summary = "Detects method expression call that can be replaced with a method call"
	info.Before = `f := foo{}
foo.bar(f)`
	info.After = `f := foo{}
f.bar()`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForExpr(&methodExprCallChecker{ctx: ctx})
	})
}

type methodExprCallChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext
}

func (c *methodExprCallChecker) VisitExpr(x ast.Expr) {
	call := astcast.ToCallExpr(x)
	s := astcast.ToSelectorExpr(call.Fun)
	id := astcast.ToIdent(s.X)

	obj := c.ctx.TypesInfo.ObjectOf(id)
	if _, ok := obj.(*types.TypeName); ok {
		c.warn(call, s)
	}
}

func (c *methodExprCallChecker) warn(cause *ast.CallExpr, s *ast.SelectorExpr) {
	selector := astcopy.SelectorExpr(s)
	selector.X = cause.Args[0]

	c.ctx.Warn(cause, "consider to change `%s` to `%s`", cause.Fun, selector)
}
