package rules

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("modeIsSubset", func() {
	It("it compares modes correctly", func() {
		Expect(modeIsSubset(0o600, 0o600)).To(BeTrue())
		Expect(modeIsSubset(0o400, 0o600)).To(BeTrue())
		Expect(modeIsSubset(0o644, 0o600)).To(BeFalse())
		Expect(modeIsSubset(0o466, 0o600)).To(BeFalse())
	})
})
