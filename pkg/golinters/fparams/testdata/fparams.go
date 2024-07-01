//golangcitest:args -Efparams
package testdata

func multiLineFuncA(
	a int,
	b string,
) (
	c bool,
	d error,
) {
	return false, nil
}

func multiLineFuncB(
	a int,
	b string,
) {
	return
}

func multiLineFuncC() (
	a bool,
	b error,
) {
	return false, nil
}

func singleLineFuncA(a int, b string) (c bool, d error) {
	return false, nil
}

func singleLineFuncB(a int) (b bool, c error) {
	return false, nil
}

func singleLineFuncC(a int, b string) (c error) {
	return nil
}

func singleLineFuncD(int, string) error {
	return nil
}

func singleLineFuncE(_ int, _ string) error {
	return nil
}

func invalidArgsFuncA(a int, // want "the parameters of the function \"invalidArgsFuncA\" should be on separate lines"
	b string) {
	return
}

func invalidArgsFuncB(a, b int, // want "the parameters of the function \"invalidArgsFuncB\" should be on separate lines"
	c string) {
	return
}

func invalidArgsFuncC(a, // want "the parameters of the function \"invalidArgsFuncC\" should be on separate lines"
	b int,
	c string,
) {
	return
}

func invalidArgsFuncD( // want "the parameters of the function \"invalidArgsFuncD\" should be on separate lines"
	a, b int,
	c string,
) {
	return
}

func invalidArgsAndResultsFuncA(a int, // want "the parameters and return values of the function \"invalidArgsAndResultsFuncA\" should be on separate lines"
	b string) (c bool,
	d error) {
	return false, nil
}

func invalidArgsAndResultsFuncB(a int, b int, // want "the parameters and return values of the function \"invalidArgsAndResultsFuncB\" should be on separate lines"
	c string) (
	d bool,
	e error) {
	return false, nil
}

func invalidResultsFuncA() (a bool, // want "the return values of the function \"invalidResultsFuncA\" should be on separate lines"
	b error) {
	return false, nil
}

func invalidResultsFuncB() ( // want "the return values of the function \"invalidResultsFuncB\" should be on separate lines"
	a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncC() ( // want "the return values of the function \"invalidResultsFuncC\" should be on separate lines"
	a bool, b error) {
	return false, nil
}

func invalidResultsFuncD() ( // want "the return values of the function \"invalidResultsFuncD\" should be on separate lines"
	a, b bool,
	c error) {
	return false, false, nil
}

func invalidResultsFuncE() (bool, bool, // want "the return values of the function \"invalidResultsFuncE\" should be on separate lines"
	error) {
	return false, false, nil
}

func invalidResultsFuncF() ( // want "the return values of the function \"invalidResultsFuncF\" should be on separate lines"
	bool,
	bool,
	error) {
	return false, false, nil
}
