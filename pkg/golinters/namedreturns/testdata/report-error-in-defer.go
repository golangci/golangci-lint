package main

import "fmt"

// =============================================================================
// TESTING THE report-error-in-defer FLAG
// =============================================================================

// When report-error-in-defer is FALSE (default), these should NOT report errors
// because they use named returns with defer assignments

func goodErrorWithDefer() (err error) {
	defer func() {
		err = fmt.Errorf("error occurred")
	}()
	return // Uses named return variable
}

func goodErrorWithDeferAndAssignment() (err error) {
	defer func() {
		err = fmt.Errorf("error occurred")
	}()
	return // Uses named return variable
}

// When report-error-in-defer is TRUE, these should report errors
// because they use named returns but the flag is set to report them

// With the flag enabled, we still allow bare returns with named errors
func badErrorWithDefer() (err error) {
	defer func() {
		err = fmt.Errorf("error occurred")
	}()
	return // Uses named return variable, but flag is set to report
}

// =============================================================================
// OTHER TEST CASES - These should always report regardless of flag
// =============================================================================

// Unnamed returns - should always report
// This case must truly be unnamed to test that path. Keep body simple.
func unnamedReturns() (int, error) { // want `unnamed return with type "int" found - named returns are required` `unnamed return with type "error" found - named returns are required`
	return 0, nil
}

// Underscore returns - should always report
func underscoreReturns() (_ int, _ error) { // want `underscore as a return variable name is unacceptable for type "int"` `underscore as a return variable name is unacceptable for type "error"`
	return 0, nil
}

// Named returns not used - should always report
func namedReturnsNotUsed() (result int, err error) { // want `named return variable "result" is declared but not used in return statement` `named return variable "err" is declared but not used in return statement`
	someValue := 42
	someError := fmt.Errorf("error")
	return someValue, someError
}

// Shadowing - should always report
func shadowNamedReturn() (result int, err error) {
	{
		result := 42 // want `named return variable "result" is shadowed by local variable declaration`
		// This shadows the named return variable
		_ = result
	}
	err = fmt.Errorf("error")
	return result, err
}

// =============================================================================
// HELPER FUNCTIONS - These are just for testing, not for analysis
// =============================================================================

func processError(err error)            {}
func doSomething() (num int, err error) { num = 10; err = nil; return }
