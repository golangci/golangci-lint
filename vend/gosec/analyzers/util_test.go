package analyzers

import (
	"go/constant"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/ssa"
)

var _ = Describe("GetConstantInt64", func() {
	It("should not panic on float constants", func() {
		// Create a float constant (simulates float64(-1))
		floatVal := constant.MakeFloat64(-1.0)
		c := &ssa.Const{Value: floatVal}

		// Should return (0, false) without panicking
		val, ok := GetConstantInt64(c)
		Expect(ok).To(BeFalse())
		Expect(val).To(Equal(int64(0)))
	})
})
