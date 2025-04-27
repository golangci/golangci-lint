//golangcitest:args -Egodot
//golangcitest:expected_exitcode 0
package p

/*
This comment won't be checked in default mode
*/

// This comment will be fixed.
func godot(a, b int) int {
	// Nothing to do here
	return a + b
}

// Foo Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// Aenean rhoncus odio enim, et pulvinar libero ultrices quis.
// Nulla at erat tellus. Maecenas id dapibus velit, ut porttitor ipsum.
func Foo() {
	// nothing to do here
}
