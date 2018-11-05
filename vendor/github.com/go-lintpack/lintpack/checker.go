package lintpack

import (
	"go/ast"
)

type checkerProto struct {
	info        *CheckerInfo
	constructor func(*Context, parameters) *Checker
}

type Checker struct {
	Info *CheckerInfo

	ctx CheckerContext

	fileWalker FileWalker

	Init func(ctx *Context)
}

// Check runs rule checker over file f.
func (c *Checker) Check(f *ast.File) []Warning {
	c.ctx.warnings = c.ctx.warnings[:0]
	c.fileWalker.WalkFile(f)
	return c.ctx.warnings
}
