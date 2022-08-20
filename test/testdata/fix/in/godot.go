//golangcitest:args -Egodot
//golangcitest:expected_exitcode 0
package p

/*
This comment won't be checked in default mode
*/

// This comment will be fixed
func godot(a, b int) int {
	// Nothing to do here
	return a + b
}
