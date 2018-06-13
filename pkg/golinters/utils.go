package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/packages"

	"github.com/golangci/golangci-lint/pkg/config"
)

func formatCode(code string, cfg *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func formatCodeBlock(code string, cfg *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("```\n%s\n```", code)
}

func getASTFilesForPkg(ctx *linter.Context, pkg *packages.Package) ([]*ast.File, *token.FileSet, error) {
	filenames := pkg.Files(ctx.Cfg.Run.AnalyzeTests)
	files := make([]*ast.File, 0, len(filenames))
	var fset *token.FileSet
	for _, filename := range filenames {
		f := ctx.ASTCache.Get(filename)
		if f == nil {
			return nil, nil, fmt.Errorf("no AST for file %s in cache: %+v", filename, *ctx.ASTCache)
		}

		if f.Err != nil {
			return nil, nil, fmt.Errorf("can't load AST for file %s: %s", f.Name, f.Err)
		}

		files = append(files, f.F)
		fset = f.Fset
	}

	return files, fset, nil
}
