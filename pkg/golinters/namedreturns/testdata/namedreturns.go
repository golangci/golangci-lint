package main

import "errors"

// =============================================================================
// GOOD EXAMPLES - These should NOT trigger any reports
// =============================================================================

// Function with no return values - this is fine
func noReturnValues() {
	// No return values, so no reports needed
}

// Function with properly named returns - this is fine
func goodNamedReturns() (result int, err error) {
	result = 42
	err = nil
	return result, err
}

// Function with bare return using named returns - this is fine
func bareReturnWithNamedReturns() (result int, err error) {
	result = 42
	err = nil
	return // Uses the named return variables
}

// Method with properly named returns - this is fine
type example struct{}

func (e *example) goodMethod() (result int, err error) {
	result = 42
	err = nil
	return result, err
}

// Function literal with properly named returns - this is fine
var goodFuncLiteral = func() (result int, err error) {
	result = 42
	err = nil
	return result, err
}

// Function with error return and defer assignment - this is fine (when flag is false)
func errorWithDeferAssignment() (err error) {
	defer func() {
		err = errors.New("error occurred")
	}()
	return // Uses the named return variable
}

// =============================================================================
// BAD EXAMPLES - These SHOULD trigger reports
// =============================================================================

// Unnamed returns - should report
func unnamedReturns() (int, error) { // want `unnamed return with type "int" found - named returns are required` `unnamed return with type "error" found - named returns are required`
	return 42, errors.New("error")
}

// Single unnamed return - should report
func singleUnnamedReturn() int { // want `unnamed return with type "int" found - named returns are required`
	return 42
}

// Underscore-named returns - should report
func underscoreReturns() (_ int, _ string) { // want `underscore as a return variable name is unacceptable for type "int"` `underscore as a return variable name is unacceptable for type "string"`
	return 42, "hello"
}

// Mixed underscore and proper names - should report on underscores
func mixedUnderscoreReturns() (_ int, result string) { // want `underscore as a return variable name is unacceptable for type "int"`
	return 42, result
}

// Named returns declared but not used in return statement - should report
func namedReturnsNotUsed() (result int, err error) { // want `named return variable "result" is declared but not used in return statement` `named return variable "err" is declared but not used in return statement`
	someValue := 42
	someError := errors.New("error")
	return someValue, someError
}

// Partial usage of named returns - should report on unused ones
func partialNamedReturnUsage() (result int, err error) { // want `named return variable "result" is declared but not used in return statement`
	err = errors.New("error")
	return 42, err
}

// Named return shadowing - should report
func shadowNamedReturn() (result int, err error) {
	{
		result := 42 // want `named return variable "result" is shadowed by local variable declaration`
		// This shadows the named return variable
		_ = result
	}
	err = errors.New("error")
	return result, err
}

// Shadowing in loops - should report
func shadowInLoop() (result int, err error) {
	for result := 0; result < 10; result++ { // want `named return variable "result" is shadowed by for loop variable` `named return variable "result" is shadowed by local variable declaration`
		// shadows named return via for-init
		_ = result
	}
	err = errors.New("error")
	return result, err
}

// Shadowing in range loops - should report
func shadowInRange() (key string, value int) {
	data := map[string]int{"a": 1, "b": 2}
	for key, value := range data { // want `named return variable "key" is shadowed by range loop variable` `named return variable "value" is shadowed by range loop variable`
		// shadows named returns via range variables
		_ = key
		_ = value
	}
	return key, value
}

// Shadowing with var declarations - should report
func shadowWithVar() (result int, err error) {
	{
		var result int = 42 // want `named return variable "result" is shadowed by local variable declaration`
		// This shadows the named return variable
		_ = result
	}
	err = errors.New("error")
	return result, err
}

// =============================================================================
// HELPER FUNCTIONS - These are just for testing, not for analysis
// =============================================================================

func processError(err error)                         {}
func doSomething() (num int, err error)              { num = 10; err = nil; return }
func multierrAppendInto(_ *error, _ error) (ok bool) { ok = false; return }
