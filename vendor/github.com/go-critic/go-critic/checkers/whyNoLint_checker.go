package checkers

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
)

func init() {
	info := lintpack.CheckerInfo{
		Name:    "whyNoLint",
		Tags:    []string{"style", "experimental"},
		Summary: "Ensures that `//nolint` comments include an explanation",
		Before:  `//nolint`,
		After:   `//nolint // reason`,
	}
	re := regexp.MustCompile(`^// *nolint(?::[^ ]+)? *(.*)$`)

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForComment(&whyNoLintChecker{
			ctx: ctx,
			re:  re,
		})
	})
}

type whyNoLintChecker struct {
	astwalk.WalkHandler

	ctx *lintpack.CheckerContext
	re  *regexp.Regexp
}

func (c whyNoLintChecker) VisitComment(cg *ast.CommentGroup) {
	if strings.HasPrefix(cg.List[0].Text, "/*") {
		return
	}
	for _, comment := range cg.List {
		sl := c.re.FindStringSubmatch(comment.Text)
		if len(sl) < 2 {
			continue
		}

		if s := sl[1]; !strings.HasPrefix(s, "//") || len(strings.TrimPrefix(s, "//")) == 0 {
			c.ctx.Warn(cg, "include an explanation for nolint directive")
			return
		}
	}
}
