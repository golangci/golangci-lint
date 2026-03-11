package cwe_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/securego/gosec/v2/cwe"
)

var _ = Describe("CWE Types", func() {
	BeforeEach(func() {
	})
	Context("when consulting cwe types", func() {
		It("it should retrieves the information and download URIs", func() {
			Expect(cwe.InformationURI).To(Equal("https://cwe.mitre.org/data/published/cwe_v4.4.pdf/"))
			Expect(cwe.DownloadURI).To(Equal("https://cwe.mitre.org/data/xml/cwec_v4.4.xml.zip"))
		})

		It("it should retrieves the weakness ID and URL", func() {
			weakness := &cwe.Weakness{ID: "798"}
			Expect(weakness).ShouldNot(BeNil())
			Expect(weakness.SprintID()).To(Equal("CWE-798"))
			Expect(weakness.SprintURL()).To(Equal("https://cwe.mitre.org/data/definitions/798.html"))
		})
	})
})
