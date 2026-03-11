package gosec_test

import (
	"go/ast"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/testutils"
)

var _ = Describe("Resolve ast node to concrete value", func() {
	Context("when attempting to resolve an ast node", func() {
		It("should successfully resolve basic literal", func() {
			var basicLiteral *ast.BasicLit
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const foo = "bar"; func main(){}`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.BasicLit); ok {
					basicLiteral = node
					return false
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(basicLiteral).ShouldNot(BeNil())
			Expect(gosec.TryResolve(basicLiteral, ctx)).Should(BeTrue())
		})

		It("should successfully resolve identifier", func() {
			var ident *ast.Ident
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; var foo string = "bar"; func main(){}`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.Ident); ok {
					ident = node
					return false
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(ident).ShouldNot(BeNil())
			Expect(gosec.TryResolve(ident, ctx)).Should(BeTrue())
		})

		It("should successfully resolve variable identifier", func() {
			var ident *ast.Ident
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; import "fmt"; func main(){ x := "test"; y := x; fmt.Println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.Ident); ok && node.Name == "y" {
					ident = node
					return false
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(ident).ShouldNot(BeNil())
			Expect(gosec.TryResolve(ident, ctx)).Should(BeTrue())
		})

		It("should successfully not resolve variable identifier with no declaration", func() {
			var ident *ast.Ident
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; import "fmt"; func main(){ x := "test"; y := x; fmt.Println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.Ident); ok && node.Name == "y" {
					ident = node
					return false
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(ident).ShouldNot(BeNil())
			ident.Obj.Decl = nil
			Expect(gosec.TryResolve(ident, ctx)).Should(BeFalse())
		})

		It("should successfully resolve assign statement", func() {
			var assign *ast.AssignStmt
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ y := x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.AssignStmt); ok {
					if id, ok := node.Lhs[0].(*ast.Ident); ok && id.Name == "y" {
						assign = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(assign).ShouldNot(BeNil())
			Expect(gosec.TryResolve(assign, ctx)).Should(BeTrue())
		})

		It("should successfully not resolve assign statement without rhs", func() {
			var assign *ast.AssignStmt
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ y := x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.AssignStmt); ok {
					if id, ok := node.Lhs[0].(*ast.Ident); ok && id.Name == "y" {
						assign = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(assign).ShouldNot(BeNil())
			assign.Rhs = []ast.Expr{}
			Expect(gosec.TryResolve(assign, ctx)).Should(BeFalse())
		})

		It("should successfully not resolve assign statement with unsolvable rhs", func() {
			var assign *ast.AssignStmt
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ y := x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.AssignStmt); ok {
					if id, ok := node.Lhs[0].(*ast.Ident); ok && id.Name == "y" {
						assign = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(assign).ShouldNot(BeNil())
			assign.Rhs = []ast.Expr{&ast.CallExpr{}}
			Expect(gosec.TryResolve(assign, ctx)).Should(BeFalse())
		})

		It("should successfully resolve a binary statement", func() {
			var target *ast.BinaryExpr
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const (x = "bar"; y = "baz"); func main(){ z := x + y; println(z) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.BinaryExpr); ok {
					target = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(target).ShouldNot(BeNil())
			Expect(gosec.TryResolve(target, ctx)).Should(BeTrue())
		})

		It("should successfully resolve value spec", func() {
			var value *ast.ValueSpec
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ var y string = x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.ValueSpec); ok {
					if len(node.Names) == 1 && node.Names[0].Name == "y" {
						value = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			Expect(gosec.TryResolve(value, ctx)).Should(BeTrue())
		})
		It("should successfully not resolve value spec without values", func() {
			var value *ast.ValueSpec
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ var y string = x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.ValueSpec); ok {
					if len(node.Names) == 1 && node.Names[0].Name == "y" {
						value = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			value.Values = []ast.Expr{}
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})

		It("should successfully not resolve value spec with unsolvable value", func() {
			var value *ast.ValueSpec
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; const x = "bar"; func main(){ var y string = x; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.ValueSpec); ok {
					if len(node.Names) == 1 && node.Names[0].Name == "y" {
						value = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			value.Values = []ast.Expr{&ast.CallExpr{}}
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})

		It("should successfully resolve composite literal", func() {
			var value *ast.CompositeLit
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; func main(){ y := []string{"value1", "value2"}; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.CompositeLit); ok {
					value = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			Expect(gosec.TryResolve(value, ctx)).Should(BeTrue())
		})

		It("should successfully not resolve composite literal without elst", func() {
			var value *ast.CompositeLit
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; func main(){ y := []string{"value1", "value2"}; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.CompositeLit); ok {
					value = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			value.Elts = []ast.Expr{}
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})

		It("should successfully not resolve composite literal with unsolvable elst", func() {
			var value *ast.CompositeLit
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; func main(){ y := []string{"value1", "value2"}; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.CompositeLit); ok {
					value = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			value.Elts = []ast.Expr{&ast.CallExpr{}}
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})

		It("should successfully not resolve call expressions", func() {
			var value *ast.CallExpr
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; func main(){ y := []string{"value1", "value2"}; println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.CallExpr); ok {
					value = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})

		It("should successfully not resolve call expressions", func() {
			var value *ast.ImportSpec
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `package main; import "fmt"; func main(){ y := []string{"value1", "value2"}; fmt.Println(y) }`)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.ImportSpec); ok {
					value = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(value).ShouldNot(BeNil())
			Expect(gosec.TryResolve(value, ctx)).Should(BeFalse())
		})
	})
})
