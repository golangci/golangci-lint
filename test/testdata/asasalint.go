//golangcitest:args -Easasalint
package testdata

import "fmt"

func getArgsLength(args ...interface{}) int {
	// this line will not report as error
	fmt.Println(args)
	return len(args)
}

func checkArgsLength(args ...interface{}) int {
	return getArgsLength(args) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
}

func someCall() {
	var a = []interface{}{1, 2, 3}
	fmt.Println(checkArgsLength(a...) == getArgsLength(a)) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
	fmt.Println(checkArgsLength(a...) == getArgsLength(a...))
}
