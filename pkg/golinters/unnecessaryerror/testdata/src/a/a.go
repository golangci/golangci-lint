package a

import (
	"a/external"
	"errors"
	"log"
	"os"
)

func ifNotNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func IfNotNil() error {
	err := ifNotNil()
	if err != nil {
		return errors.New("")
	}
	return nil
}

func ifNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func IfNil() error {
	err := ifNil()
	if err == nil {
		return errors.New("")
	}
	return nil
}

func multiRet() (int, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return 1, errors.New("")
}

func MultiRet() (int, error) {
	i, err := multiRet()
	if err != nil {
		return 0, errors.New("")
	}
	return i, nil
}

func switchNotNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func SwitchNotNil() error {
	err := switchNotNil()
	switch {
	case err != nil:
		return errors.New("")
	}
	return nil
}

func multiFuncReturnErr() error {
	return errors.New("")
}

func MultiFuncReturnErr_1() error {
	err := multiFuncReturnErr()
	if err != nil {
		return nil
	}
	return nil
}

func MultiFuncReturnErr_2() error {
	err := multiFuncReturnErr()
	if err != nil {
		return err // OK: returned from exported function.
	}
	return nil
}

func multiFuncReturnMultiErr() (error, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New(""), errors.New("")
}

func MultiFuncReturnMultiErr() error {
	err1, err2 := multiFuncReturnMultiErr()
	if err1 != nil {
		return err2 // err2 OK: returned from exported function.
	}
	return nil
}

func namedReturn() (err error) { // want "error is only ever nil-checked; consider returning a bool instead"
	err = errors.New("")
	return
}

func NamedReturn() {
	err := namedReturn()
	if err != nil {
		return
	}
}

type Receiver1 struct{}

func (r *Receiver1) methodReturnsErr() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func (r *Receiver1) MethodReturnsErr() {
	err := r.methodReturnsErr()
	if err != nil {
		return
	}
}

func passedToExternal() error {
	return errors.New("")
}

func PassedToExternal() {
	err := passedToExternal()
	log.Print(err) // OK: passed to external function.
}

func passedToBuiltinFunc() error {
	return errors.New("")
}

func PassedToBuiltinFunc() {
	err := passedToBuiltinFunc()
	println(err) // OK: passed to builtin function.
}

func funcPassedToExternal_1() error {
	return errors.New("")
}

func funcPassedToExternal_2() {
	err := external.FuncPassedToExternal(funcPassedToExternal_1) // OK: passed to external function.
	log.Print(err)
}

func funcPassedToExternalMethod_1() error {
	return errors.New("")
}

func funcPassedToExternalMethod_2() {
	err := (&external.Exported{}).FuncPassedToExternalMethod(funcPassedToExternalMethod_1) // OK: passed to external function.
	log.Print(err)
}

func methodCalled() error {
	return errors.New("")
}

func MethodCalled() string {
	err := methodCalled()
	return err.Error() // OK: method called.
}

func typeAssert() error {
	return errors.New("")
}

func TypeAssert() {
	err := typeAssert()
	if _, ok := err.(*os.SyscallError); ok { // OK: type asserted.
		return
	}
}

func typeSwitch() error {
	return errors.New("")
}

func TypeSwitch() {
	err := typeSwitch()
	switch err.(type) { // OK: type switch.
	case *os.PathError:
		return
	}
}

func returnInStruct() error {
	return errors.New("")
}

type _returnInStruct struct {
	err error
}

func ReturnInStruct() _returnInStruct {
	err := returnInStruct()
	if err != nil {
		return _returnInStruct{err: err} // OK: returned from exported function.
	}
	return _returnInStruct{}
}

func assignedToVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func AssignedToVar() {
	err := assignedToVar()
	x := err
	_ = x
}

func structField() error {
	return errors.New("")
}

func StructField() {
	x := struct{ err error }{err: structField()} // OK: assigned to struct field.
	_ = x
}

func mapIndex() error {
	return errors.New("")
}

func MapIndex() {
	m := map[int]error{}
	m[0] = mapIndex() // OK: inserted into map.
}

func sliceIndex() error {
	return errors.New("")
}

func SliceIndex() {
	s := make([]error, 1)
	s[0] = sliceIndex() // OK: inserted into slice.
}

func Closure() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if err != nil {
		return
	}
	return
}

func assignedToExportedGlobalVar1() error {
	return errors.New("")
}

var AssignedToExportedGlobalVar1 = assignedToExportedGlobalVar1() // OK: assigned to global variable.

func assignedToExportedGlobalVar2() error {
	return errors.New("")
}

var AssignedToExportedGlobalVar2_ error

func AssignedToExportedGlobalVar2() {
	AssignedToExportedGlobalVar2_ = assignedToExportedGlobalVar2()
}

func assignedToUnexportedGlobalVar() error {
	return errors.New("")
}

var _assignedToUnexportedGlobalVar = assignedToUnexportedGlobalVar() // OK: assigned to global variable.

func genericFunc1[T any]() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericFunc1_1() {
	err := genericFunc1[int]()
	if err != nil {
		return
	}
}

func GenericFunc1_2() {
	err := genericFunc1[string]()
	if err != nil {
		return
	}
}

func genericFunc2[T any]() error {
	return errors.New("")
}

func GenericFunc2_1() {
	err := genericFunc2[int]()
	if err != nil {
		return
	}
}

func GenericFunc2_2() error {
	err := genericFunc2[string]()
	if err != nil {
		return err // OK: returned from exported function.
	}
	return nil
}

func genericError[T any]() T {
	var zero T
	return zero
}

func GenericError() {
	if err := genericError[error](); err != nil { // OK: ignore generic return types.
		return
	}
	return
}

func genericError2[T error]() T {
	var zero T
	return zero
}

func GenericError2() {
	if err := genericError2[error](); err != nil { // OK: ignore generic return types.
		return
	}
	return
}

func callchain1() error {
	return errors.New("")
}

func callchain1_1(err error) error {
	return err
}

func callchain1_2(err error) error {
	return err
}

func Callchain1() error {
	err := callchain1_2(callchain1_1(callchain1()))
	if err != nil {
		return err
	}
	return nil
}

func callchain2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain2_1(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return err
}

func callchain2_2(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return err
}

func Callchain2() error {
	err := callchain2_2(callchain2_1(callchain2()))
	if err != nil {
		return nil
	}
	return nil
}

func callchain3() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain3_1() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain3()
}

func callchain3_2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain3_1()
}
func Callchain3() {
	err := callchain3_2()
	if err != nil {
		return
	}
}

// functions starting with an underscore are considered unexported.
func _specialUnexported() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func SpecialUnexported() {
	err := _specialUnexported()
	if err != nil {
		return
	}
}

func nilCheckInDefer() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckInDefer() {
	err := nilCheckInDefer()
	defer func() {
		if err != nil {
			return
		}
	}()
}

func nilCheckInForLoop() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckInForLoop() {
	err := nilCheckInForLoop()
	for err != nil {
		return
	}
}

func closurePassedToExternal() error {
	return errors.New("")
}

func ClosurePassedToExternal() {
	err := external.ClosurePassedToExternal(func() error { // OK: passed to external function.
		return closurePassedToExternal()
	})
	log.Print(err)
}

func assignedInClosureToExternal() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func AssignedInClosureToExternal() {
	var err1 error
	err2 := external.AssignedInClosure(func() error {
		err1 = assignedInClosureToExternal()
		return errors.New("")
	})
	log.Print(err1, err2)
}

func assignedToExternalStructField1() error {
	return errors.New("")
}

func AssignedToExternalStructField1() {
	s := external.AssignedToExternalStructField1{
		Err: assignedToExternalStructField1(),
	}
	_ = s
}

func assignedToExternalStructField2() error {
	return errors.New("")
}

func assignedToExternalStructField2_1() error {
	err := assignedToExternalStructField2()
	if err != nil {
		return err
	}
	return nil
}

func AssignedToExternalStructField2() {
	s := external.AssignedToExternalStructField2{
		Err: assignedToExternalStructField2_1(), // OK: Assigned to struct field.
	}
	_ = s
}

func funcAssignedToExternalStructField1() error {
	return errors.New("")
}

func FuncAssignedToExternalStructField1() {
	s := external.FuncAssignedToExternalStructField1{
		Func: funcAssignedToExternalStructField1, // OK: Func assigned to struct field.
	}
	_ = s
}

func funcAssignedToExternalStructField2() error {
	return errors.New("")
}

func funcAssignedToExternalStructField2_1() error {
	err := funcAssignedToExternalStructField2()
	if err != nil {
		return err
	}
	return nil
}

func FuncAssignedToExternalStructField2() {
	s := external.FuncAssignedToExternalStructField2{
		Func: funcAssignedToExternalStructField2_1, // OK: Caller assigned to struct field.
	}
	_ = s
}

type _methodAssignedToExternalStructField struct{}

func (s _methodAssignedToExternalStructField) methodAssignedToExternalStructField() error {
	return errors.New("")
}

func (s _methodAssignedToExternalStructField) methodAssignedToExternalStructField1() error {
	err := s.methodAssignedToExternalStructField()
	if err != nil {
		return err
	}
	return nil
}

func MethodAssignedToExternalStructField() {
	var x _methodAssignedToExternalStructField
	s := external.MethodAssignedToExternalStructField{
		Func: x.methodAssignedToExternalStructField1,
	}
	_ = s
}
