//golangcitest:args -Egodot
package testdata

// want +2 "Comment should end in a period"

// Godot checks top-level comments
func Godot() {
	// nothing to do here
}

// want +4 "Comment should end in a period"

// Foo Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// Aenean rhoncus odio enim, et pulvinar libero ultrices quis.
// Nulla at erat tellus. Maecenas id dapibus velit, ut porttitor ipsum
func Foo() {
	// nothing to do here
}
