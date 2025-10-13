//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package splitseq

import (
	"bytes"
	"strings"
)

func _() {
	for _, line := range strings.Split("", "") { // want "Ranging over SplitSeq is more efficient"
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
	for _ = range strings.Split("", "") { // want "Ranging over SplitSeq is more efficient"
	}
	for range strings.Split("", "") { // want "Ranging over SplitSeq is more efficient"
	}
	for range bytes.Split(nil, nil) { // want "Ranging over SplitSeq is more efficient"
	}
	{
		lines := strings.Split("", "") // want "Ranging over SplitSeq is more efficient"
		for _, line := range lines {
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
