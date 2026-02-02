//golangcitest:args -Enosharedparamtype
//golangcitest:expected_exitcode 0
package testdata

type MyType struct{}

func funcWithSharedParamTypes(x int, y int) int {
	return x + y
}

func funcWithThreeSharedParams(a string, b string, c string) string {
	return a + b + c
}

func funcWithMixedParams(a string, b int, c int) string {
	return a
}

func funcWithSharedPointers(x *int, y *int) int {
	return *x + *y
}

func funcWithSharedSlices(x []string, y []string) int {
	return len(x) + len(y)
}

func (m MyType) methodWithSharedParams(x int, y int) int {
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
