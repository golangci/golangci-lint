//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package fieldsseq

import (
	"bytes"
	"strings"
)

func _() {
	for line := range strings.FieldsSeq("") { // want "Ranging over FieldsSeq is more efficient"
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
	for range strings.FieldsSeq("") { // want "Ranging over FieldsSeq is more efficient"
	}
	for range strings.FieldsSeq("") { // want "Ranging over FieldsSeq is more efficient"
	}
	for range bytes.FieldsSeq(nil) { // want "Ranging over FieldsSeq is more efficient"
	}
	{
		lines := strings.FieldsSeq("") // want "Ranging over FieldsSeq is more efficient"
		for line := range lines {
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
