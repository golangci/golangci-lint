//golangcitest:args -Egocheckdirectives
package testdata

// Okay cases:

//go:generate echo hello world

//go:embed
var Value string

//go:

// Problematic cases:

// go:embed // want "go directive contains leading space: // go:embed"

//    go:embed // want "go directive contains leading space: //    go:embed"

//go:genrate // want "unrecognized go directive: //go:genrate"
