package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var HandlerAnalyzer = &analysis.Analyzer{
	Name:     "extractorlint",
	Doc:      "Checks and validates ElementHandler defs.",
	Run:      runHandlerAnalyzer,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runHandlerAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(node ast.Node) {
		genDecl := node.(*ast.GenDecl)

		hl := parseHandler(genDecl)
		if hl == nil {
			return
		}

		issues := hl.validate()
		for _, issue := range issues {
			pass.Report(issue.Diagnose())
		}
	})

	return nil, nil //nolint:nilnil
}
