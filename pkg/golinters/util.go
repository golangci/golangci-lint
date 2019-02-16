package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	gopackages "golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func formatCode(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func formatCodeBlock(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("```\n%s\n```", code)
}

func getAllFileNames(ctx *linter.Context) []string {
	var ret []string
	uniqFiles := map[string]bool{} // files are duplicated for test packages
	for _, pkg := range ctx.Packages {
		for _, f := range pkg.GoFiles {
			if uniqFiles[f] {
				continue
			}
			uniqFiles[f] = true
			ret = append(ret, f)
		}
	}
	return ret
}

func getASTFilesForGoPkg(ctx *linter.Context, pkg *gopackages.Package) ([]*ast.File, *token.FileSet, error) {
	var files []*ast.File
	var fset *token.FileSet
	for _, filename := range pkg.GoFiles {
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
