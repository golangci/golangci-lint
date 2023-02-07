//golangcitest:args -Elll
//golangcitest:config_path testdata/configs/lll.yml
package testdata

import (
	_ "embed"
	_ "unsafe"
)

func Lll() {
	// In my experience, long lines are the lines with comments, not the code. So this is a long comment // want "line is 137 characters"
}

//go:generate mockgen -source lll.go -destination a_verylong_generate_mock_my_lll_interface.go --package testdata -self_package github.com/golangci/golangci-lint/test/testdata
type MyLllInterface interface {
}

//go:linkname VeryLongNameForTestAndLinkNameFunction github.com/golangci/golangci-lint/test/testdata.VeryLongNameForTestAndLinkedNameFunction
func VeryLongNameForTestAndLinkNameFunction()

func VeryLongNameForTestAndLinkedNameFunction() {}

//go:embed this_is_a_text_file_with_a_very_long_name_repeat_this_is_a_text_file_with_a_very_long_name_repeat_this_is_a_text_file_with_a_very_long_name.txt
var _ string
