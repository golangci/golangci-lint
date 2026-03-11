package testutils

import (
	"go/ast"

	"github.com/securego/gosec/v2"
)

// MockVisitor is useful for stubbing out ast.Visitor with callback
// and looking for specific conditions to exist.
type MockVisitor struct {
	Context  *gosec.Context
	Callback func(n ast.Node, ctx *gosec.Context) bool
}

// NewMockVisitor creates a new empty struct, the Context and
// Callback must be set manually. See call_list_test.go for an example.
func NewMockVisitor() *MockVisitor {
	return &MockVisitor{}
}

// Visit satisfies the ast.Visitor interface
func (v *MockVisitor) Visit(n ast.Node) ast.Visitor {
	if v.Callback(n, v.Context) {
		return v
	}
	return nil
}
