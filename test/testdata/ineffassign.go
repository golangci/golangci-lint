//args: -Eineffassign
package testdata

func _() {
	x := 0
	for {
		_ = x
		x = 0 // ERROR "ineffectual assignment to `x`"
		x = 0
	}
}
