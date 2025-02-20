package p_test

import "fmt"

func WithGolintIssues(b bool) { //nolint:staticcheck
	if b {
		return
	} else {
		fmt.Print("1")
	}
}
