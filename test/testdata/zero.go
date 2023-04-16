//golangcitest:args -Ezero
package testdata

var _ string = "" // want "shoud not assign zero value"

func _() {
	n := 0 // want "shoud not assign zero value"
	_ = n

	var _ []int = nil     // want "shoud not assign zero value"
	var _ []int = []int{} // OK
	m := int32(0)         // OK
	_ = m
	var _ *int = nil            // want "shoud not assign zero value"
	var _ struct{} = struct{}{} // want "shoud not assign zero value"
	var _, _ int                // OK
	var _, _ int = 0, 1         // want "shoud not assign zero value"
	var _, _ int = 1, 2         // OK
	var _, _ int = 1 - 1, 2 - 2 // want "shoud not assign zero value"
	var _ bool = false          // want "shoud not assign zero value"
	var _ bool = true           // OK

	type T struct{ N int }
	var _ T = T{} // want "shoud not assign zero value"

	{
		n, _ := func() (int, int) { return 0, 0 }() // OK
		_ = n
	}

	{
		var n, _ = func() (int, int) { return 0, 0 }() // OK
		_ = n
	}
}
