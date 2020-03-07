package checks

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	config "github.com/tommy-muehle/go-mnd/config"
)

const AssignCheck = "assign"

type AssignAnalyzer struct {
	pass   *analysis.Pass
	config *config.Config
}

func NewAssignAnalyzer(pass *analysis.Pass, config *config.Config) *AssignAnalyzer {
	return &AssignAnalyzer{
		pass:   pass,
		config: config,
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
		if a.isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, AssignCheck)
		}
	case *ast.BinaryExpr:
		a.checkBinaryExpr(x)
	}
}

func (a *AssignAnalyzer) checkBinaryExpr(expr *ast.BinaryExpr) {
	switch x := expr.X.(type) {
	case *ast.BasicLit:
		if a.isMagicNumber(x) {
			a.pass.Reportf(x.Pos(), reportMsg, x.Value, AssignCheck)
		}
	}

	switch y := expr.Y.(type) {
	case *ast.BasicLit:
		if a.isMagicNumber(y) {
			a.pass.Reportf(y.Pos(), reportMsg, y.Value, AssignCheck)
		}
	}
}

func (a *AssignAnalyzer) isMagicNumber(l *ast.BasicLit) bool {
	return (l.Kind == token.FLOAT || l.Kind == token.INT) && !a.config.IsIgnoredNumber(l.Value)
}
