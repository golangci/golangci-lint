//golangcitest:args -Efparams
//golangcitest:expected_exitcode 0
package fparams

func invalidArgsFuncA(a int,
	b string) {
	return
}

func invalidArgsFuncB(a, b int,
	c string) {
	return
}

func invalidArgsFuncC(a,
	b int,
	c string,
) {
	return
}

func invalidArgsFuncD(
	a, b int,
	c string,
) {
	return
}

func invalidArgsAndResultsFuncA(a int,
	b string) (c bool,
	d error) {
	return false, nil
}

func invalidArgsAndResultsFuncB(a int, b int,
	c string) (
	d bool,
	e error) {
	return false, nil
}

func invalidResultsFuncA() (a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncB() (
	a bool,
	b error) {
	return false, nil
}

func invalidResultsFuncC() (
	a bool, b error) {
	return false, nil
}

func invalidResultsFuncD() (
	a, b bool,
	c error) {
	return false, false, nil
}

func invalidResultsFuncE() (bool, bool,
	error) {
	return false, false, nil
}

func invalidResultsFuncF() (
	bool,
	bool,
	error) {
	return false, false, nil
}
