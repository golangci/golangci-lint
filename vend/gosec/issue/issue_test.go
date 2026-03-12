package issue_test

import (
	"go/ast"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/rules"
	"github.com/securego/gosec/v2/testutils"
)

var _ = Describe("Issue", func() {
	Context("when creating a new issue", func() {
		It("should create a code snippet from the specified ast.Node", func() {
			var target *ast.BasicLit
			source := `package main
			const foo = "bar"
			func main(){
				println(foo)
			}
			`
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", source)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.BasicLit); ok {
					target = node
					return false
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(target).ShouldNot(BeNil())

			fobj := ctx.GetFileAtNodePos(target)
			issue := issue.New(fobj, target, "TEST", "", issue.High, issue.High)
			Expect(issue).ShouldNot(BeNil())
			Expect(issue.Code).Should(MatchRegexp(`"bar"`))
			Expect(issue.Line).Should(Equal("2"))
			Expect(issue.Col).Should(Equal("16"))
			Expect(issue.Cwe).Should(BeNil())
		})

		It("should return an error if specific context is not able to be obtained", func() {
			Skip("Not implemented")
		})

		It("should construct file path based on line and file information", func() {
			var target *ast.AssignStmt

			source := `package main
			import "fmt"
			func main() {
				username := "admin"
				password := "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
				fmt.Println("Doing something with: ", username, password)
			}`

			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", source)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.AssignStmt); ok {
					if id, ok := node.Lhs[0].(*ast.Ident); ok && id.Name == "password" {
						target = node
					}
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(target).ShouldNot(BeNil())

			// Use hardcoded rule to check assignment
			cfg := gosec.NewConfig()
			rule, _ := rules.NewHardcodedCredentials("TEST", cfg)
			foundIssue, err := rule.Match(target, ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(foundIssue).ShouldNot(BeNil())
			Expect(foundIssue.FileLocation()).Should(MatchRegexp("foo.go:5"))
		})

		It("should provide accurate line and file information", func() {
			Skip("Not implemented")
		})

		It("should provide accurate line and file information for multi-line statements", func() {
			var target *ast.CallExpr
			source := `
package main
import (
   	"net"
)
func main() {
	_, _ := net.Listen("tcp", 
	"0.0.0.0:2000")
}
`
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", source)
			ctx := pkg.CreateContext("foo.go")
			v := testutils.NewMockVisitor()
			v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
				if node, ok := n.(*ast.CallExpr); ok {
					target = node
				}
				return true
			}
			v.Context = ctx
			ast.Walk(v, ctx.Root)
			Expect(target).ShouldNot(BeNil())

			cfg := gosec.NewConfig()
			rule, _ := rules.NewBindsToAllNetworkInterfaces("TEST", cfg)
			issue, err := rule.Match(target, ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(issue).ShouldNot(BeNil())
			Expect(issue.File).Should(MatchRegexp("foo.go"))
			Expect(issue.Line).Should(MatchRegexp("7-8"))
			Expect(issue.Col).Should(Equal("10"))
		})

		It("should maintain the provided severity score", func() {
			Skip("Not implemented")
		})

		It("should maintain the provided confidence score", func() {
			Skip("Not implemented")
		})
	})
})
