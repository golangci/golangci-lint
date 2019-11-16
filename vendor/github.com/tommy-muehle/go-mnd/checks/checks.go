package checks

import (
	"go/ast"
	"go/token"
)

const reportMsg = "Magic number: %v, in <%s> detected"

func isMagicNumber(l *ast.BasicLit) bool {
	return (l.Kind == token.FLOAT || l.Kind == token.INT) && l.Value != "0"
}
