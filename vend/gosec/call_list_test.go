package gosec_test

import (
	"go/ast"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/testutils"
)

var _ = Describe("Call List", func() {
	var calls gosec.CallList
	BeforeEach(func() {
		calls = gosec.NewCallList()
	})

	It("should not return any matches when empty", func() {
		Expect(calls.Contains("foo", "bar")).Should(BeFalse())
	})

	It("should be possible to add a single call", func() {
		Expect(calls).Should(BeEmpty())
		calls.Add("foo", "bar")
		Expect(calls).Should(HaveLen(1))

		expected := make(map[string]bool)
		expected["bar"] = true
		actual := map[string]bool(calls["foo"])
		Expect(actual).Should(Equal(expected))
	})

	It("should be possible to add multiple calls at once", func() {
		Expect(calls).Should(BeEmpty())
		calls.AddAll("fmt", "Sprint", "Sprintf", "Printf", "Println")

		expected := map[string]bool{
			"Sprint":  true,
			"Sprintf": true,
			"Printf":  true,
			"Println": true,
		}
		actual := map[string]bool(calls["fmt"])
		Expect(actual).Should(Equal(expected))
	})

	It("should be possible to add pointer call", func() {
		Expect(calls).Should(BeEmpty())
		calls.Add("*bytes.Buffer", "WriteString")
		actual := calls.ContainsPointer("*bytes.Buffer", "WriteString")
		Expect(actual).Should(BeTrue())
	})

	It("should be possible to check pointer call", func() {
		Expect(calls).Should(BeEmpty())
		calls.Add("bytes.Buffer", "WriteString")
		actual := calls.ContainsPointer("*bytes.Buffer", "WriteString")
		Expect(actual).Should(BeTrue())
	})

	It("should not return a match if none are present", func() {
		calls.Add("ioutil", "Copy")
		Expect(calls.Contains("fmt", "Println")).Should(BeFalse())
	})

	It("should match a call based on selector and ident", func() {
		calls.Add("ioutil", "Copy")
		Expect(calls.Contains("ioutil", "Copy")).Should(BeTrue())
	})

	It("should match a package call expression", func() {
		// Create file to be scanned
		pkg := testutils.NewTestPackage()
		defer pkg.Close()
		pkg.AddFile("md5.go", testutils.SampleCodeG401[0].Code[0])

		ctx := pkg.CreateContext("md5.go")

		// Search for md5.New()
		calls.Add("crypto/md5", "New")

		// Stub out visitor and count number of matched call expr
		matched := 0
		v := testutils.NewMockVisitor()
		v.Context = ctx
		v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
			if _, ok := n.(*ast.CallExpr); ok && calls.ContainsPkgCallExpr(n, ctx, false) != nil {
				matched++
			}
			return true
		}
		ast.Walk(v, ctx.Root)
		Expect(matched).Should(Equal(1))
	})

	It("should match a package call expression", func() {
		// Create file to be scanned
		pkg := testutils.NewTestPackage()
		defer pkg.Close()
		pkg.AddFile("cipher.go", testutils.SampleCodeG405[0].Code[0])

		ctx := pkg.CreateContext("cipher.go")

		// Search for des.NewCipher()
		calls.Add("crypto/des", "NewCipher")

		// Stub out visitor and count number of matched call expr
		matched := 0
		v := testutils.NewMockVisitor()
		v.Context = ctx
		v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
			if _, ok := n.(*ast.CallExpr); ok && calls.ContainsPkgCallExpr(n, ctx, false) != nil {
				matched++
			}
			return true
		}
		ast.Walk(v, ctx.Root)
		Expect(matched).Should(Equal(1))
	})

	It("should match a package call expression", func() {
		// Create file to be scanned
		pkg := testutils.NewTestPackage()
		defer pkg.Close()
		pkg.AddFile("md4.go", testutils.SampleCodeG406[0].Code[0])

		ctx := pkg.CreateContext("md4.go")

		// Search for md4.New()
		calls.Add("golang.org/x/crypto/md4", "New")

		// Stub out visitor and count number of matched call expr
		matched := 0
		v := testutils.NewMockVisitor()
		v.Context = ctx
		v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
			if _, ok := n.(*ast.CallExpr); ok && calls.ContainsPkgCallExpr(n, ctx, false) != nil {
				matched++
			}
			return true
		}
		ast.Walk(v, ctx.Root)
		Expect(matched).Should(Equal(1))
	})

	It("should match a call expression", func() {
		// Create file to be scanned
		pkg := testutils.NewTestPackage()
		defer pkg.Close()
		pkg.AddFile("main.go", testutils.SampleCodeG104[6].Code[0])

		ctx := pkg.CreateContext("main.go")

		calls.Add("bytes.Buffer", "WriteString")
		calls.Add("strings.Builder", "WriteString")
		calls.Add("io.Pipe", "CloseWithError")
		calls.Add("fmt", "Fprintln")

		// Stub out visitor and count number of matched call expr
		matched := 0
		v := testutils.NewMockVisitor()
		v.Context = ctx
		v.Callback = func(n ast.Node, ctx *gosec.Context) bool {
			if _, ok := n.(*ast.CallExpr); ok && calls.ContainsCallExpr(n, ctx) != nil {
				matched++
			}
			return true
		}
		ast.Walk(v, ctx.Root)
		Expect(matched).Should(Equal(5))
	})
})
