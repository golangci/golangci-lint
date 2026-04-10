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

	v1 = 1
	v2 = value()
	v3 = struct{}{}
	v5 = StructT{}

	v3 = struct{}{}
	v4 = struct{}{}

	v1 = 1
	_ = 2

	_ = v1
	_ = v2
	_ = v3
	_ = v4
	_ = v5
}

const cpkgi1 = 1
const cpkgi2 = 2

const (
	cpkgig1 = 1
	cpkgig2 = 2
)

const cpkgit1 int = 1
const cpkgit2 int = 2

const (
	cpkgitg1 int = 1
	cpkgitg2 int = 2
)

func var_func_init() {
	var (
		vpkgi1 = 1
		vpkgi2 = 2
	)

	var (
		vpkgi5 = struct{}{}
		vpkgi6 = struct{}{}
	)

	var (
		vpkgi7 = StructT{}
		vpkgi8 = StructT{}
	)

	var (
		vpkgi9  = 1
		vpkgi10 = value()
	)

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
	v1 int
	v2 int
	v3 string `tag:"value"`
	v4 string `tag:"value"`
}
