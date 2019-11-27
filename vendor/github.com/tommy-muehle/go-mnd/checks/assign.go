package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const AssignCheck = "assign"

type AssignAnalyzer struct {
	pass *analysis.Pass
}

func NewAssignAnalyzer(pass *analysis.Pass) *AssignAnalyzer {
	return &AssignAnalyzer{
		pass: pass,
	}
}

func (a *AssignAnalyzer) NodeFilter() []ast.Node {
	return []ast.Node{
		(*ast.KeyValueExpr)(nil),
	}
}

func (a *AssignAnalyzer) Check(n ast.Node) {
	expr, ok := n.(*ast.KeyValueExpr)
	if !ok {
		return
	}

	switch x := expr.Value.(type) {
	case *ast.BasicLit:
		if isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, AssignCheck)
		}
	case *ast.BinaryExpr:
		a.checkBinaryExpr(x)
	}
}

func (a *AssignAnalyzer) checkBinaryExpr(expr *ast.BinaryExpr) {
	switch x := expr.X.(type) {
	case *ast.BasicLit:
		if isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, AssignCheck)
		}
	}

	switch y := expr.Y.(type) {
	case *ast.BasicLit:
		if isMagicNumber(y) {
			a.pass.Reportf(y.Pos(), reportMsg, y.Value, AssignCheck)
		}
	}
}
