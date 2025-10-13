//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package stringscutprefix

import (
	"strings"
)

var (
	s, pre, suf string
)

// test supported cases of pattern 1 - CutPrefix
func _() {
	if after, ok := strings.CutPrefix(s, pre); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		a := after
		_ = a
	}
	if after, ok := strings.CutPrefix("", ""); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		a := after
		_ = a
	}
	if after, ok := strings.CutPrefix(s, ""); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		println([]byte(after))
	}
	if after, ok := strings.CutPrefix(s, ""); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		a, b := "", after
		_, _ = a, b
	}
	if after, ok := strings.CutPrefix(s, ""); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		a, b := after, strings.TrimPrefix(s, "") // only replace the first occurrence
		s = "123"
		b = strings.TrimPrefix(s, "") // only replace the first occurrence
		_, _ = a, b
	}

	var a, b string
	if after, ok := strings.CutPrefix(s, ""); ok { // want "HasPrefix \\+ TrimPrefix can be simplified to CutPrefix"
		a, b = "", after
		_, _ = a, b
	}
}

// test basic cases for CutSuffix - only covering the key differences
func _() {
	if before, ok := strings.CutSuffix(s, suf); ok { // want "HasSuffix \\+ TrimSuffix can be simplified to CutSuffix"
		a := before
		_ = a
	}
	if before, ok := strings.CutSuffix(s, ""); ok { // want "HasSuffix \\+ TrimSuffix can be simplified to CutSuffix"
		println([]byte(before))
	}
}

// test cases that are not supported by pattern1 - CutPrefix
func _() {
	ok := strings.HasPrefix("", "")
	if ok { // noop, currently it doesn't track the result usage of HasPrefix
		a := strings.TrimPrefix("", "")
		_ = a
	}
	if strings.HasPrefix(s, pre) {
		a := strings.TrimPrefix("", "") // noop, as the argument isn't the same
		_ = a
	}
	if strings.HasPrefix(s, pre) {
		var result string
		result = strings.TrimPrefix("", "") // noop, as we believe define is more popular.
		_ = result
	}
	if strings.HasPrefix("", "") {
		a := strings.TrimPrefix(s, pre) // noop, as the argument isn't the same
		_ = a
	}
	if s1 := s; strings.HasPrefix(s1, pre) {
		a := strings.TrimPrefix(s1, pre) // noop, as IfStmt.Init is present
		_ = a
	}
}

// test basic unsupported case for CutSuffix
func _() {
	if strings.HasSuffix(s, suf) {
		a := strings.TrimSuffix("", "") // noop, as the argument isn't the same
		_ = a
	}
}

var value0 string

// test supported cases of pattern2 - CutPrefix
func _() {
	if after, ok := strings.CutPrefix(s, pre); ok { // want "TrimPrefix can be simplified to CutPrefix"
		println(after)
	}
	if after, ok := strings.CutPrefix(s, pre); ok { // want "TrimPrefix can be simplified to CutPrefix"
		println(after)
	}
	if after, ok := strings.CutPrefix(s, pre); ok { // want "TrimPrefix can be simplified to CutPrefix"
		println(strings.TrimPrefix(s, pre)) // noop here
	}
	if after, ok := strings.CutPrefix(s, ""); ok { // want "TrimPrefix can be simplified to CutPrefix"
		println(after)
	}
	var ok bool // define an ok variable to test the fix won't shadow it for its if stmt body
	_ = ok
	if after, ok0 := strings.CutPrefix(s, pre); ok0 { // want "TrimPrefix can be simplified to CutPrefix"
		println(after)
	}
	var predefined string
	if predefined = strings.TrimPrefix(s, pre); s != predefined { // noop
		println(predefined)
	}
	if predefined = strings.TrimPrefix(s, pre); s != predefined { // noop
		println(&predefined)
	}
	var value string
	if value = strings.TrimPrefix(s, pre); s != value { // noop
		println(value)
	}
	lhsMap := make(map[string]string)
	if lhsMap[""] = strings.TrimPrefix(s, pre); s != lhsMap[""] { // noop
		println(lhsMap[""])
	}
	arr := make([]string, 0)
	if arr[0] = strings.TrimPrefix(s, pre); s != arr[0] { // noop
		println(arr[0])
	}
	type example struct {
		field string
	}
	var e example
	if e.field = strings.TrimPrefix(s, pre); s != e.field { // noop
		println(e.field)
	}
}

// test basic cases for pattern2 - CutSuffix
func _() {
	if before, ok := strings.CutSuffix(s, suf); ok { // want "TrimSuffix can be simplified to CutSuffix"
		println(before)
	}
	if before, ok := strings.CutSuffix(s, suf); ok { // want "TrimSuffix can be simplified to CutSuffix"
		println(before)
	}
}

// test cases that not supported by pattern2 - CutPrefix
func _() {
	if after := strings.TrimPrefix(s, pre); s != pre { // noop
		println(after)
	}
	if after := strings.TrimPrefix(s, pre); after != pre { // noop
		println(after)
	}
	if strings.TrimPrefix(s, pre) != s {
		println(strings.TrimPrefix(s, pre))
	}
}

// test basic unsupported case for pattern2 - CutSuffix
func _() {
	if before := strings.TrimSuffix(s, suf); s != suf { // noop
		println(before)
	}
}
