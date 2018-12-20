package checkers

import (
	"go/ast"
	"go/constant"
	"strings"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
	"github.com/go-toolsmith/astcast"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "flagName"
	info.Tags = []string{"diagnostic", "experimental"}
	info.Summary = "Detects flag names with whitespace"
	info.Before = `b := flag.Bool(" foo ", false, "description")`
	info.After = `b := flag.Bool("foo", false, "description")`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForExpr(&flagNameChecker{ctx: ctx})
	})
}

type flagNameChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext
}

func (c *flagNameChecker) VisitExpr(expr ast.Expr) {
	call := astcast.ToCallExpr(expr)
	switch qualifiedName(call.Fun) {
	case "flag.Bool", "flag.Duration", "flag.Float64", "flag.String",
		"flag.Int", "flag.Int64", "flag.Uint", "flag.Uint64":
		c.checkFlagName(call, call.Args[0])
	case "flag.BoolVar", "flag.DurationVar", "flag.Float64Var", "flag.StringVar",
		"flag.IntVar", "flag.Int64Var", "flag.UintVar", "flag.Uint64Var":
		c.checkFlagName(call, call.Args[1])
	}
}

func (c *flagNameChecker) checkFlagName(call *ast.CallExpr, arg ast.Expr) {
	cv := c.ctx.TypesInfo.Types[arg].Value
	if cv == nil {
		return // Non-constant name
	}
	name := constant.StringVal(cv)
	if strings.Contains(name, " ") {
		c.warnWhitespace(call, name)
	}
}

func (c *flagNameChecker) warnWhitespace(cause ast.Node, name string) {
	c.ctx.Warn(cause, "flag name %q contains whitespace", name)
}
