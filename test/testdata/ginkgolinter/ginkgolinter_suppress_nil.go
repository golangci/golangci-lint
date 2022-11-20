//golangcitest:args --disable-all -Eginkgolinter
//golangcitest:config_path configs/ginkgolinter_suppress_nil.yml
package ginkgolinter

import (
	"errors"
	. "github.com/onsi/gomega"
)

func LenUsecase_nil() {
	var fakeVarUnderTest []int
	Expect(fakeVarUnderTest).Should(BeEmpty())     // valid
	Expect(fakeVarUnderTest).ShouldNot(HaveLen(5)) // valid

	Expect(len(fakeVarUnderTest)).Should(Equal(0))           // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.Should\\(BeEmpty\\(\\)\\). instead"
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(2))        // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.ShouldNot\\(HaveLen\\(2\\)\\). instead"
	Expect(len(fakeVarUnderTest)).To(BeNumerically("==", 0)) // // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.To\\(BeEmpty\\(\\)\\). instead"

	fakeVarUnderTest = append(fakeVarUnderTest, 3)
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(0))        // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.ShouldNot\\(BeEmpty\\(\\)\\). instead"
	Expect(len(fakeVarUnderTest)).Should(Equal(1))           // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.Should\\(HaveLen\\(1\\)\\). instead"
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">", 0))  // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\). instead"
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">=", 1)) // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\). instead"
	Expect(len(fakeVarUnderTest)).To(BeNumerically("!=", 0)) // want "ginkgo-linter: wrong length assertion; consider using .Expect\\(fakeVarUnderTest\\)\\.ToNot\\(BeEmpty\\(\\)\\). instead"
}

func NilUsecase_nil() {
	y := 5
	x := &y
	Expect(x == nil).To(Equal(true))
	Expect(nil == x).To(Equal(true))
	Expect(x != nil).To(Equal(true))
	Expect(x == nil).To(BeTrue())
	Expect(x == nil).To(BeFalse())
}
func BooleanUsecase_nil() {
	x := true
	Expect(x).To(Equal(true)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeTrue\\(\\)\\). instead"
	x = false
	Expect(x).To(Equal(false)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeFalse\\(\\)\\). instead"
}

func ErrorUsecase_nil() {
	err := errors.New("fake error")
	funcReturnsErr := func() error { return err }

	Expect(err).To(BeNil())              // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\). instead"
	Expect(err == nil).To(Equal(true))   // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\). instead"
	Expect(err == nil).To(BeFalse())     // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\). instead"
	Expect(err != nil).To(BeTrue())      // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\). instead"
	Expect(funcReturnsErr()).To(BeNil()) // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(funcReturnsErr\\(\\)\\)\\.To\\(Succeed\\(\\)\\). instead"
}
