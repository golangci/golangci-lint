package cwe_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2/cwe"
)

var _ = Describe("CWE data", func() {
	BeforeEach(func() {
	})
	Context("when consulting cwe data", func() {
		It("it should retrieves the weakness", func() {
			weakness := cwe.Get("798")
			Expect(weakness).ShouldNot(BeNil())
			Expect(weakness.ID).ShouldNot(BeNil())
			Expect(weakness.Name).ShouldNot(BeNil())
			Expect(weakness.Description).ShouldNot(BeNil())
		})
	})
})
