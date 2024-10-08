package main

import "fmt"

func main() {
	if false {
		fmt.Println("Nested Hello, world") // Comment inside if.
		panic("this is bad")               /* want "panic call without same line comment justifying it" */
	} else {
		panic("this is ok") // A comment on the same line makes panic() ok.
	}
}
