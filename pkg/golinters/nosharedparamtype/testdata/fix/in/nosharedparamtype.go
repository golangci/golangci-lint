//golangcitest:args -Enosharedparamtype
//golangcitest:expected_exitcode 0
package testdata

type MyType struct{}

func funcWithSharedParamTypes(x, y int) int {
	return x + y
}

func funcWithThreeSharedParams(a, b, c string) string {
	return a + b + c
}

func funcWithMixedParams(a string, b, c int) string {
	return a
}

func funcWithSharedPointers(x, y *int) int {
	return *x + *y
}

func funcWithSharedSlices(x, y []string) int {
	return len(x) + len(y)
}

func (m MyType) methodWithSharedParams(x, y int) int {
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

func (m MyType) methodWithSeparateParams(x float64, y int) int {
	return int(x) + y
}
