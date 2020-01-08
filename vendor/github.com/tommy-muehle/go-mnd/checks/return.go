package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const ReturnCheck = "return"

type ReturnAnalyzer struct {
	pass *analysis.Pass
}

func NewReturnAnalyzer(pass *analysis.Pass) *ReturnAnalyzer {
	return &ReturnAnalyzer{
		pass: pass,
	}
}

func (a *ReturnAnalyzer) NodeFilter() []ast.Node {
	return []ast.Node{
		(*ast.ReturnStmt)(nil),
	}
}

func (a *ReturnAnalyzer) Check(n ast.Node) {
	stmt, ok := n.(*ast.ReturnStmt)
	if !ok {
		return
	}

	for _, expr := range stmt.Results {
		switch x := expr.(type) {
		case *ast.BasicLit:
			if isMagicNumber(x) {
				a.pass.Reportf(x.Pos(), reportMsg, x.Value, ReturnCheck)
			}
		}
	}
}
