//golangcitest:config_path testdata/ginkgolinter_suppress_async.yml
//golangcitest:args -Eginkgolinter
package ginkgolinter

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func LenUsecase_async() {
	var fakeVarUnderTest []int
	Expect(fakeVarUnderTest).Should(BeEmpty())     // valid
	Expect(fakeVarUnderTest).ShouldNot(HaveLen(5)) // valid

	Expect(len(fakeVarUnderTest)).Should(Equal(0))           // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.Should\\(BeEmpty\\(\\)\\). instead"
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(2))        // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.ShouldNot\\(HaveLen\\(2\\)\\)"
	Expect(len(fakeVarUnderTest)).To(BeNumerically("==", 0)) // // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.To\\(BeEmpty\\(\\)\\)"

	fakeVarUnderTest = append(fakeVarUnderTest, 3)
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(0))        // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.ShouldNot\\(BeEmpty\\(\\)\\)"
	Expect(len(fakeVarUnderTest)).Should(Equal(1))           // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.Should\\(HaveLen\\(1\\)\\)"
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">", 0))  // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\)"
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">=", 1)) // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\)"
	Expect(len(fakeVarUnderTest)).To(BeNumerically("!=", 0)) // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\)"
}

func NilUsecase_async() {
	y := 5
	x := &y
	Expect(x == nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion. Consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\)"
	Expect(nil == x).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion. Consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\)"
	Expect(x != nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion. Consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\)"
	Expect(x == nil).To(BeTrue())    // want "ginkgo-linter: wrong nil assertion. Consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\)"
	Expect(x == nil).To(BeFalse())   // want "ginkgo-linter: wrong nil assertion. Consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\)"
}
func BooleanUsecase_async() {
	x := true
	Expect(x).To(Equal(true)) // want "ginkgo-linter: wrong boolean assertion. Consider using .Expect\\(x\\)\\.To\\(BeTrue\\(\\)\\)"
	x = false
	Expect(x).To(Equal(false)) // want "ginkgo-linter: wrong boolean assertion. Consider using .Expect\\(x\\)\\.To\\(BeFalse\\(\\)\\)"
}

func ErrorUsecase_async() {
	err := errors.New("fake error")
	funcReturnsErr := func() error { return err }

	Expect(err).To(BeNil())              // want "ginkgo-linter: wrong error assertion. Consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\)"
	Expect(err == nil).To(Equal(true))   // want "ginkgo-linter: wrong error assertion. Consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\)"
	Expect(err == nil).To(BeFalse())     // want "ginkgo-linter: wrong error assertion. Consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\)"
	Expect(err != nil).To(BeTrue())      // want "ginkgo-linter: wrong error assertion. Consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\)"
	Expect(funcReturnsErr()).To(BeNil()) // want "ginkgo-linter: wrong error assertion. Consider using .Expect\\(funcReturnsErr\\(\\)\\)\\.To\\(Succeed\\(\\)\\)"
}

func HaveLen0Usecase_async() {
	x := make([]string, 0)
	Expect(x).To(HaveLen(0)) // want "ginkgo-linter: wrong length assertion. Consider using .Expect\\(x\\)\\.To\\(BeEmpty\\(\\)\\)"
}

func WrongComparisonUsecase_async() {
	x := 8
	Expect(x == 8).To(BeTrue())    // want "ginkgo-linter: wrong comparison assertion. Consider using .Expect\\(x\\)\\.To\\(Equal\\(8\\)\\)"
	Expect(x < 9).To(BeTrue())     // want "ginkgo-linter: wrong comparison assertion. Consider using .Expect\\(x\\)\\.To\\(BeNumerically\\(\"<\", 9\\)\\)"
	Expect(x < 7).To(Equal(false)) // want "ginkgo-linter: wrong comparison assertion. Consider using .Expect\\(x\\)\\.ToNot\\(BeNumerically\\(\"<\", 7\\)\\)"

	p1, p2 := &x, &x
	Expect(p1 == p2).To(Equal(true)) // want "ginkgo-linter: wrong comparison assertion. Consider using .Expect\\(p1\\).To\\(BeIdenticalTo\\(p2\\)\\)"
}

func slowInt_async() int {
	time.Sleep(time.Second)
	return 42
}

// WrongEventuallyWithFunction_async Should trigger no warning
func WrongEventuallyWithFunction_async() {
	Eventually(slowInt_async).Should(Equal(42))   // valid
	Eventually(slowInt_async()).Should(Equal(42)) // suppressed
}

var _ = FDescribe("Should warn for focused containers", func() {
	FContext("should not allow FContext", func() {
		FWhen("should not allow FWhen", func() {
			FIt("should not allow FIt", func() {})
		})
	})
})
