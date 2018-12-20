package checkers

import (
	"go/ast"
	"regexp"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "docStub"
	info.Tags = []string{"style", "experimental"}
	info.Summary = "Detects comments that silence go lint complaints about doc-comment"
	info.Before = `
// Foo ...
func Foo() {
}`
	info.After = `
// (A) - remove the doc-comment stub
func Foo() {}
// (B) - replace it with meaningful comment
// Foo is a demonstration-only function.
func Foo() {}`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		c := &docStubChecker{ctx: ctx}
		c.badCommentRE = regexp.MustCompile(`//\s?\w+([^a-zA-Z]+|( XXX.?))$`)
		return astwalk.WalkerForFuncDecl(c)
	})
}

type docStubChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext

	badCommentRE *regexp.Regexp
}

func (c *docStubChecker) VisitFuncDecl(decl *ast.FuncDecl) {
	if decl.Name.IsExported() && decl.Doc != nil && c.badCommentRE.MatchString(decl.Doc.List[0].Text) {
		c.warn(decl)
	}
}

func (c *docStubChecker) warn(decl *ast.FuncDecl) {
	c.ctx.Warn(decl, "silencing go lint doc-comment warnings is unadvised")
}
