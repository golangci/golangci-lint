//golangcitest:args -Emultisplit
//golangcitest:config_path testdata/multisplit.yml
package multisplit

import (
	_ "errors"
	_ "strings"
)

func value() int {
	return 1
}

func value2() (int, int) {
	return 1, 2
}

type StructT struct{}

func assign() {
	var (
		v1 int
		v2 int
		v3 struct{}
		v4 struct{}
		v5 StructT
	)

	v1, v2, v3, v5 = 1, value(), struct{}{}, StructT{} // want `assignment with multiple targets \(v1, v2, v3, v5\) should be split into individual assignments`

	v3, v4 = struct{}{}, struct{}{} // want `assignment with multiple targets \(v3, v4\) should be split into individual assignments`

	v1, v2 = 1, 2 // want `assignment with multiple targets \(v1, v2\) should be split into individual assignments`
	v1, _ = 1, 2 // want `assignment with multiple targets \(v1, _\) should be split into individual assignments`

	_ = v1
	_ = v2
	_ = v3
	_ = v4
	_ = v5
}

const cpkgi1, cpkgi2 = 1, 2 // want `const declaration with multiple identifiers \(cpkgi1, cpkgi2\) should be split into individual declarations`

const (
	cpkgig1, cpkgig2 = 1, 2 // want `const declaration with multiple identifiers \(cpkgig1, cpkgig2\) should be split into individual declarations`
)

const cpkgit1, cpkgit2 int = 1, 2 // want `const declaration with multiple identifiers \(cpkgit1, cpkgit2\) should be split into individual declarations`

const (
	cpkgitg1, cpkgitg2 int = 1, 2 // want `const declaration with multiple identifiers \(cpkgitg1, cpkgitg2\) should be split into individual declarations`
)

func var_func_init() {
	var vpkgi1, vpkgi2 = 1, 2 // want `variable declaration with multiple identifiers and initializers \(vpkgi1, vpkgi2\) should be split into individual declarations`

	var vpkgi3, vpkgi4 = 3, 4 // want `variable declaration with multiple identifiers and initializers \(vpkgi3, vpkgi4\) should be split into individual declarations`

	var vpkgi5, vpkgi6 = struct{}{}, struct{}{} // want `variable declaration with multiple identifiers and initializers \(vpkgi5, vpkgi6\) should be split into individual declarations`

	var vpkgi7, vpkgi8 = StructT{}, StructT{} // want `variable declaration with multiple identifiers and initializers \(vpkgi7, vpkgi8\) should be split into individual declarations`

	var vpkgi9, vpkgi10 = 1, value() // want `variable declaration with multiple identifiers and initializers \(vpkgi9, vpkgi10\) should be split into individual declarations`

	var vpkgi11, vpkgi12 = value2()

	_ = vpkgi1
	_ = vpkgi2
	_ = vpkgi3
	_ = vpkgi4
	_ = vpkgi5
	_ = vpkgi6
	_ = vpkgi7
	_ = vpkgi8
	_ = vpkgi9
	_ = vpkgi10
	_ = vpkgi11
	_ = vpkgi12
}

type s1 struct {
	v1, v2 int // want `struct field declaration with multiple identifiers \(v1, v2\) should be split into individual fields`
	v3, v4 string `tag:"value"` // want `struct field declaration with multiple identifiers \(v3, v4\) should be split into individual fields`
}
