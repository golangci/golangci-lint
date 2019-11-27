package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const ConditionCheck = "condition"

type ConditionAnalyzer struct {
	pass *analysis.Pass
}

func NewConditionAnalyzer(pass *analysis.Pass) *ConditionAnalyzer {
	return &ConditionAnalyzer{
		pass: pass,
	}
}

func (a *ConditionAnalyzer) NodeFilter() []ast.Node {
	return []ast.Node{
		(*ast.IfStmt)(nil),
	}
}

func (a *ConditionAnalyzer) Check(n ast.Node) {
	expr, ok := n.(*ast.IfStmt).Cond.(*ast.BinaryExpr)
	if !ok {
		return
	}

	switch x := expr.X.(type) {
	case *ast.BasicLit:
		if isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, ConditionCheck)
		}
	}

	switch y := expr.Y.(type) {
	case *ast.BasicLit:
		if isMagicNumber(y) {
			a.pass.Reportf(y.Pos(), reportMsg, y.Value, ConditionCheck)
		}
	}
}
