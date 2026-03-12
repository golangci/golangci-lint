package gosec_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/testutils"
)

var _ = Describe("Import Tracker", func() {
	Context("when tracking a file", func() {
		It("should parse the imports from file", func() {
			tracker := gosec.NewImportTracker()
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package foo
				import "fmt"
				func foo() {
				  fmt.Println()
				}
			`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			pkgs := pkg.Pkgs()
			Expect(pkgs).Should(HaveLen(1))
			files := pkgs[0].Syntax
			Expect(files).Should(HaveLen(1))
			tracker.TrackFile(files[0])
			Expect(tracker.Imported).Should(Equal(map[string][]string{"fmt": {"fmt"}}))
		})
		It("should parse the named imports from file", func() {
			tracker := gosec.NewImportTracker()
			pkg := testutils.NewTestPackage()
			defer pkg.Close()
			pkg.AddFile("foo.go", `
				package foo
				import fm "fmt"
				func foo() {
				  fm.Println()
				}
			`)
			err := pkg.Build()
			Expect(err).ShouldNot(HaveOccurred())
			pkgs := pkg.Pkgs()
			Expect(pkgs).Should(HaveLen(1))
			files := pkgs[0].Syntax
			Expect(files).Should(HaveLen(1))
			tracker.TrackFile(files[0])
			Expect(tracker.Imported).Should(Equal(map[string][]string{"fmt": {"fm"}}))
		})
	})
})
