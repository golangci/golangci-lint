//golangcitest:args -Elll
//golangcitest:config_path testdata/lll.yml
package testdata

import (
	_ "unsafe"
)

func Lll() {
	// want +1 "line is 141 characters"
	// In my experience, long lines are the lines with comments, not the code. So this is a long comment, a very long comment, yes very long.
}

//go:generate mockgen -source lll.go -destination a_verylong_generate_mock_my_lll_interface.go --package testdata -self_package github.com/golangci/golangci-lint/test/testdata
type MyLllInterface interface {
}

//go:linkname VeryLongNameForTestAndLinkNameFunction github.com/golangci/golangci-lint/test/testdata.VeryLongNameForTestAndLinkedNameFunction
func VeryLongNameForTestAndLinkNameFunction()

func VeryLongNameForTestAndLinkedNameFunction() {}
