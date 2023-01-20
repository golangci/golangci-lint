//golangcitest:config_path configs/ginkgolinter_suppress_err.yml
//golangcitest:args --disable-all -Eginkgolinter
package ginkgolinter

import (
	"errors"
	. "github.com/onsi/gomega"
)

func LenUsecase_err() {
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

func NilUsecase_err() {
	y := 5
	x := &y
	Expect(x == nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(nil == x).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(x != nil).To(Equal(true)) // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\). instead"
	Expect(x == nil).To(BeTrue())    // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.To\\(BeNil\\(\\)\\). instead"
	Expect(x == nil).To(BeFalse())   // want "ginkgo-linter: wrong nil assertion; consider using .Expect\\(x\\)\\.ToNot\\(BeNil\\(\\)\\). instead"
}
func BooleanUsecase_err() {
	x := true
	Expect(x).To(Equal(true)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeTrue\\(\\)\\). instead"
	x = false
	Expect(x).To(Equal(false)) // want "ginkgo-linter: wrong boolean assertion; consider using .Expect\\(x\\)\\.To\\(BeFalse\\(\\)\\). instead"
}

func ErrorUsecase_err() {
	err := errors.New("fake error")
	funcReturnsErr := func() error { return err }

	Expect(err).To(BeNil())
	Expect(err == nil).To(Equal(true))
	Expect(err == nil).To(BeFalse())
	Expect(err != nil).To(BeTrue())
	Expect(funcReturnsErr()).To(BeNil())
}
