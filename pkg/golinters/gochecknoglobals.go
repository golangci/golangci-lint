package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gochecknoglobals struct{}

func (Gochecknoglobals) Name() string {
	return "gochecknoglobals"
}

func (Gochecknoglobals) Desc() string {
	return "Checks that no globals are present in Go code"
}

func (lint Gochecknoglobals) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		res = append(res, lint.checkFile(f.F, f.Fset)...)
	}

	return res, nil
}

func (lint Gochecknoglobals) checkFile(f *ast.File, fset *token.FileSet) []result.Issue {
	var res []result.Issue
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for _, vn := range valueSpec.Names {
				if isWhitelisted(vn) {
					continue
				}

				res = append(res, result.Issue{
					Pos:        fset.Position(vn.Pos()),
					Text:       fmt.Sprintf("%s is a global variable", formatCode(vn.Name, nil)),
					FromLinter: lint.Name(),
				})
			}
		}
	}

	return res
}

type whitelistedExpression struct {
	Name    string
	SelName string
}

func isWhitelisted(v ast.Node) bool {
	switch i := v.(type) {
	case *ast.Ident:
		return i.Name == "_" || i.Name == "version" || looksLikeError(i)
	case *ast.CallExpr:
		if expr, ok := i.Fun.(*ast.SelectorExpr); ok {
			return isWhitelistedSelectorExpression(expr)
		}
	case *ast.CompositeLit:
		if expr, ok := i.Type.(*ast.SelectorExpr); ok {
			return isWhitelistedSelectorExpression(expr)
		}
	}

	return false
}

func isWhitelistedSelectorExpression(v *ast.SelectorExpr) bool {
	x, ok := v.X.(*ast.Ident)
	if !ok {
		return false
	}

	whitelist := []whitelistedExpression{
		{
			Name:    "errors",
			SelName: "New",
		},
		{
			Name:    "fmt",
			SelName: "Errorf",
		},
		{
			Name:    "regexp",
			SelName: "MustCompile",
		},
	}

	for _, i := range whitelist {
		if x.Name == i.Name && v.Sel.Name == i.SelName {
			return true
		}
	}

	return false
}

// looksLikeError returns true if the AST identifier starts
// with 'err' or 'Err', or false otherwise.
//
// TODO: https://github.com/leighmcculloch/gochecknoglobals/issues/5
func looksLikeError(i *ast.Ident) bool {
	prefix := "err"
	if i.IsExported() {
		prefix = "Err"
	}
	return strings.HasPrefix(i.Name, prefix)
}
