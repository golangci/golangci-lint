package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "input file")
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file %s: %s", filename, err)
	}
	ast.Print(fset, f)
}
