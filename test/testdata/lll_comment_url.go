//golangcitest:args -Elll
//golangcitest:config_path testdata/configs/lll.yml
package testdata

import (
	"fmt"
	_ "unsafe"
)

// https://github.com/golangci/golangci-lint/blob/master/pkg/result/processors/testdata/autogen_exclude_block_comment.go
func lllCommentURL1() {
	// https://github.com/golangci/golangci-lint/blob/master/pkg/result/processors/testdata/autogen_exclude_block_comment.go
	fmt.Println("lll")
}

// https://github.com/golangci/golangci-lint/blob/master/pkg/result/processors/testdata/autogen_exclude_block_comment.go foobar // want "line is 160 characters"
func lllCommentURL2() {
	// https://github.com/golangci/golangci-lint/blob/master/pkg/result/processors/testdata/autogen_exclude_block_comment.go foobar // want "line is 164 characters"
	fmt.Println("lll")
}
