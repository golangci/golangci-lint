//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package splitseq

import (
	"bytes"
	"strings"
)

func _() {
	for line := range strings.SplitSeq("", "") { // want "Ranging over SplitSeq is more efficient"
		println(line)
	}
	for i, line := range strings.Split("", "") { // nope: uses index var
		println(i, line)
	}
	for i, _ := range strings.Split("", "") { // nope: uses index var
		println(i)
	}
	for i := range strings.Split("", "") { // nope: uses index var
		println(i)
	}
	for range strings.SplitSeq("", "") { // want "Ranging over SplitSeq is more efficient"
	}
	for range strings.SplitSeq("", "") { // want "Ranging over SplitSeq is more efficient"
	}
	for range bytes.SplitSeq(nil, nil) { // want "Ranging over SplitSeq is more efficient"
	}
	{
		lines := strings.SplitSeq("", "") // want "Ranging over SplitSeq is more efficient"
		for line := range lines {
			println(line)
		}
	}
	{
		lines := strings.Split("", "") // nope: lines is used not just by range
		for _, line := range lines {
			println(line)
		}
		println(lines)
	}
}
