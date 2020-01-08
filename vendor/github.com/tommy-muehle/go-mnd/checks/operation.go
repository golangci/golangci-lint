package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const OperationCheck = "operation"

type OperationAnalyzer struct {
	pass *analysis.Pass
}

func NewOperationAnalyzer(pass *analysis.Pass) *OperationAnalyzer {
	return &OperationAnalyzer{
		pass: pass,
	}
}

func (a *OperationAnalyzer) NodeFilter() []ast.Node {
	return []ast.Node{
		(*ast.AssignStmt)(nil),
	}
}

func (a *OperationAnalyzer) Check(n ast.Node) {
	stmt, ok := n.(*ast.AssignStmt)
	if !ok {
		return
	}

	for _, expr := range stmt.Rhs {
		switch x := expr.(type) {
		case *ast.BinaryExpr:
			switch xExpr := x.X.(type) {
			case *ast.BinaryExpr:
				a.checkBinaryExpr(xExpr)
			}
			switch yExpr := x.Y.(type) {
			case *ast.BinaryExpr:
				a.checkBinaryExpr(yExpr)
			}

			a.checkBinaryExpr(x)
		}
	}
}

func (a *OperationAnalyzer) checkBinaryExpr(expr *ast.BinaryExpr) {
	switch x := expr.X.(type) {
	case *ast.BasicLit:
		if isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, OperationCheck)
		}
	}

	switch y := expr.Y.(type) {
	case *ast.BasicLit:
		if isMagicNumber(y) {
			a.pass.Reportf(y.Pos(), reportMsg, y.Value, OperationCheck)
		}
	}
}
