//golangcitest:args -Emultisplit
//golangcitest:config_path testdata/multisplit.yml
//golangcitest:expected_exitcode 0
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

	v1, v2, v3, v5 = 1, value(), struct{}{}, StructT{}

	v3, v4 = struct{}{}, struct{}{}

	v1, _ = 1, 2

	_ = v1
	_ = v2
	_ = v3
	_ = v4
	_ = v5
}

const cpkgi1, cpkgi2 = 1, 2

const (
	cpkgig1, cpkgig2 = 1, 2
)

const cpkgit1, cpkgit2 int = 1, 2

const (
	cpkgitg1, cpkgitg2 int = 1, 2
)

func var_func_init() {
	var vpkgi1, vpkgi2 = 1, 2

	var vpkgi5, vpkgi6 = struct{}{}, struct{}{}

	var vpkgi7, vpkgi8 = StructT{}, StructT{}

	var vpkgi9, vpkgi10 = 1, value()

	var vpkgi11, vpkgi12 = value2()

	_ = vpkgi1
	_ = vpkgi2
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
	v1, v2 int
	v3, v4 string `tag:"value"`
}
