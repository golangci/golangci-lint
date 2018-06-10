package p_test

import "fmt"

func WithGolintIssues(b bool) { //nolint:megacheck
	if b {
		return
	} else {
		fmt.Print("1")
	}
}
