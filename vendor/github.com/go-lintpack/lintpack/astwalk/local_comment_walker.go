package astwalk

import (
	"go/ast"
	"strings"
)

type localCommentWalker struct {
	visitor LocalCommentVisitor
}

func (w *localCommentWalker) WalkFile(f *ast.File) {
	if !w.visitor.EnterFile(f) {
		return
	}

	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if !ok || !w.visitor.EnterFunc(decl) {
			continue
		}

		for _, cg := range f.Comments {
			// Not sure that decls/comments are sorted
			// by positions, so do a naive full scan for now.
			if cg.Pos() < decl.Pos() || cg.Pos() > decl.End() {
				continue
			}

			var group []*ast.Comment
			visitGroup := func(list []*ast.Comment) {
				if len(list) == 0 {
					return
				}
				cg := &ast.CommentGroup{List: list}
				w.visitor.VisitLocalComment(cg)
			}
			for _, comment := range cg.List {
				if strings.HasPrefix(comment.Text, "/*") {
					visitGroup(group)
					group = group[:0]
					visitGroup([]*ast.Comment{comment})
				} else {
					group = append(group, comment)
				}
			}
			visitGroup(group)
		}
	}
}
