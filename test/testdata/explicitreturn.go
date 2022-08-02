//golangcitest:args -Eexplicitreturn
package testdata

// This is the sort of thing we want to discourage because this will implicitly return x and y
func SimpleNakedReturn() (x int, y string) {
	return // ERROR `implicit return found in function SimpleNakedReturn`
}

// If your function doesn't return anything, the linter is happy to let you keep your naked return
// ...even if it is redundant here
func IgnorableProcedure() {
	return
}

// Or even no return statement at all
func AnotherIgnorableProcedure() {}

// If you simply return the expected values in your return statement, everyone's happy
func SimpleFunction() (x int, y string) {
	return x, y // Return values are as foretold in the function signature
}

// If you have multiple return paths in your function, the linter will only catch the offensive ones
func MultipleReturnPaths() (x int, y string) {
	if true {
		return x, y // All good here, but...
	}
	return // ERROR `implicit return found in function MultipleReturnPaths`
}

// Also works on anonymous functions!
var x = func() (x int, y string) { return } // ERROR `implicit return found in anonymous function`

// But we don't require the number of elements in the return statement to match the number in the function signature
func NestedGoodAnonymousFunction() (x int, y string) {
	var f = func() (x int, y string) { return x, y } // This is fine
	return f()                                       // Your compiler will catch most of the stupid stuff you could try here
}
