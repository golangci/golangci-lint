//golangcitest:args -Enosharedparamtype
package testdata

type MyType struct{}

func funcWithSharedParamTypes(x, y int) int { // want "function funcWithSharedParamTypes has parameters using shared type"
	return x + y
}

func funcWithThreeSharedParams(a, b, c string) string { // want "function funcWithThreeSharedParams has parameters using shared type"
	return a + b + c
}

func funcWithMixedParams(a string, b, c int) string { // want "function funcWithMixedParams has parameters using shared type"
	return a
}

func funcWithSharedPointers(x, y *int) int { // want "function funcWithSharedPointers has parameters using shared type"
	return *x + *y
}

func funcWithSharedSlices(x, y []string) int { // want "function funcWithSharedSlices has parameters using shared type"
	return len(x) + len(y)
}

func (m MyType) methodWithSharedParams(x, y int) int { // want "function methodWithSharedParams has parameters using shared type"
	return x + y
}

func funcWithoutSharedParamTypes(x int, y int) int {
	return x + y
}

func funcWithSingleParam(x int) int {
	return x
}

func funcWithNoParams() int {
	return 42
}

func (m MyType) methodWithSeparateParams(x int, y int) int {
	return x + y
}
