package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/tools/imports"
)

const newPkgName = "two"

const (
	srcDir = "./pkg/config"
	dstDir = "./pkg/commands/internal/migrate/two"
)

func main() {
	stat, err := os.Stat(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		log.Fatalf("%s is not a directory", srcDir)
	}

	_ = os.RemoveAll(dstDir)

	err = processPackage(srcDir, dstDir)
	if err != nil {
		log.Fatalf("Processing package error: %v", err)
	}
}

func processPackage(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(srcPath string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if skipFile(srcPath) {
			return nil
		}

		fset := token.NewFileSet()

		file, err := parser.ParseFile(fset, srcPath, nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", srcPath, err)
		}

		processFile(file)

		return writeNewFile(fset, file, srcPath, dstDir)
	})
}

func skipFile(path string) bool {
	if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
		return true
	}

	if filepath.Base(path) == "base_loader.go" {
		return true
	}

	if filepath.Base(path) == "loader.go" {
		return true
	}

	return false
}

func processFile(file *ast.File) {
	file.Name.Name = newPkgName

	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			continue

		case *ast.GenDecl:
			switch d.Tok {
			case token.CONST, token.VAR:
				continue
			case token.TYPE:
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					processStructFields(structType)
				}
			default:
				// noop
			}

			newDecls = append(newDecls, decl)
		}
	}

	file.Decls = newDecls
}

func processStructFields(structType *ast.StructType) {
	var newFields []*ast.Field

	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 && !field.Names[0].IsExported() {
			continue
		}

		if field.Tag == nil {
			continue
		}

		field.Type = convertType(field.Type)
		field.Tag.Value = convertStructTag(field.Tag.Value)

		newFields = append(newFields, field)
	}

	structType.Fields.List = newFields
}

func convertType(expr ast.Expr) ast.Expr {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return expr
	}

	switch ident.Name {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "float32", "float64":
		return &ast.StarExpr{X: ident}

	default:
		return expr
	}
}

func convertStructTag(value string) string {
	structTag := reflect.StructTag(strings.Trim(value, "`"))

	key := structTag.Get("mapstructure")

	if key == ",squash" {
		return wrapStructTag(`yaml:",inline"`)
	}

	return wrapStructTag(fmt.Sprintf(`yaml:"%[1]s,omitempty" toml:"%[1]s,omitempty"`, key))
}

func wrapStructTag(s string) string {
	return "`" + s + "`"
}

func writeNewFile(fset *token.FileSet, file *ast.File, srcPath, dstDir string) error {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("printing %s: %w", srcPath, err)
	}

	dstPath := filepath.Join(dstDir, filepath.Base(srcPath))

	_ = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)

	formatted, err := imports.Process(dstPath, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("formatting %s: %w", dstPath, err)
	}

	//nolint:gosec,mnd // The permission is right.
	err = os.WriteFile(dstPath, formatted, 0o644)
	if err != nil {
		return fmt.Errorf("writing file %s: %w", dstPath, err)
	}

	return nil
}
