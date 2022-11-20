//golangcitest:config_path configs/ginkgolinter_suppress_len.yml
//golangcitest:args --disable-all -Eginkgolinter
package ginkgolinter

import (
	"errors"
	. "github.com/onsi/gomega"
)

func LenUsecase_len() {
	var fakeVarUnderTest []int
	Expect(fakeVarUnderTest).Should(BeEmpty())     // valid
	Expect(fakeVarUnderTest).ShouldNot(HaveLen(5)) // valid

	Expect(len(fakeVarUnderTest)).Should(Equal(0))
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(2))
	Expect(len(fakeVarUnderTest)).To(BeNumerically("==", 0))

	fakeVarUnderTest = append(fakeVarUnderTest, 3)
	Expect(len(fakeVarUnderTest)).ShouldNot(Equal(0))
	Expect(len(fakeVarUnderTest)).Should(Equal(1))
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">", 0))
	Expect(len(fakeVarUnderTest)).To(BeNumerically(">=", 1))
	Expect(len(fakeVarUnderTest)).To(BeNumerically("!=", 0))
}

func NilUsecase_len() {
	y := 5
	x := &y
	Expect(x == nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(nil == x).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(x != nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\). instead"
	Expect(x == nil).To(BeTrue())    // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(x == nil).To(BeFalse())   // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\). instead"
}
func BooleanUsecase_len() {
	x := true
	Expect(x).To(Equal(true)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeTrue\\(\\)\\). instead"
	x = false
	Expect(x).To(Equal(false)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeFalse\\(\\)\\). instead"
}

func ErrorUsecase_len() {
	err := errors.New("fake error")
	funcReturnsErr := func() error { return err }

	Expect(err).To(BeNil())              // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\). instead"
	Expect(err == nil).To(Equal(true))   // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.ToNot\\(HaveOccurred\\(\\)\\). instead"
	Expect(err == nil).To(BeFalse())     // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\). instead"
	Expect(err != nil).To(BeTrue())      // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(err\\)\\.To\\(HaveOccurred\\(\\)\\). instead"
	Expect(funcReturnsErr()).To(BeNil()) // want "ginkgo-linter: wrong error assertion; consider using .Expect\\(funcReturnsErr\\(\\)\\)\\.To\\(Succeed\\(\\)\\). instead"
}
