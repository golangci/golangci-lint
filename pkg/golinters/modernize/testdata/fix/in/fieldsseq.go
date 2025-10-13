//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package fieldsseq

import (
	"bytes"
	"strings"
)

func _() {
	for _, line := range strings.Fields("") { // want "Ranging over FieldsSeq is more efficient"
		println(line)
	}
	for i, line := range strings.Fields("") { // nope: uses index var
		println(i, line)
	}
	for i, _ := range strings.Fields("") { // nope: uses index var
		println(i)
	}
	for i := range strings.Fields("") { // nope: uses index var
		println(i)
	}
	for _ = range strings.Fields("") { // want "Ranging over FieldsSeq is more efficient"
	}
	for range strings.Fields("") { // want "Ranging over FieldsSeq is more efficient"
	}
	for range bytes.Fields(nil) { // want "Ranging over FieldsSeq is more efficient"
	}
	{
		lines := strings.Fields("") // want "Ranging over FieldsSeq is more efficient"
		for _, line := range lines {
			println(line)
		}
	}
	{
		lines := strings.Fields("") // nope: lines is used not just by range
		for _, line := range lines {
			println(line)
		}
		println(lines)
	}
}
