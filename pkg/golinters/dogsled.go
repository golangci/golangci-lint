package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Dogsled struct{}

func (Dogsled) Name() string {
	return "dogsled"
}

func (Dogsled) Desc() string {
	return "Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())"
}

func (d Dogsled) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {

	var res []result.Issue
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		v := returnsVisitor{
			maxBlanks: lintCtx.Settings().Dogsled.MaxBlankIdentifiers,
			f:         f.Fset,
		}
		ast.Walk(&v, f.F)
		res = append(res, v.issues...)
	}

	return res, nil
}

type returnsVisitor struct {
	f         *token.FileSet
	maxBlanks int
	issues    []result.Issue
}

func (v *returnsVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}
	if funcDecl.Body == nil {
		return v
	}

	for _, expr := range funcDecl.Body.List {
		assgnStmt, ok := expr.(*ast.AssignStmt)
		if !ok {
			continue
		}

		numBlank := 0
		for _, left := range assgnStmt.Lhs {
			ident, ok := left.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == "_" {
				numBlank++
			}
		}

		if numBlank > v.maxBlanks {
			v.issues = append(v.issues, result.Issue{
				FromLinter: Dogsled{}.Name(),
				Text:       fmt.Sprintf("declaration has %v blank identifiers", numBlank),
				Pos:        v.f.Position(assgnStmt.Pos()),
			})
		}
	}
	return v
}
