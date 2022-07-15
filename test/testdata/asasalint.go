//golangcitest:args -Easasalint
package testdata

import "fmt"

func getArgsLength(args ...interface{}) int {
	return len(args)
}

func checkArgsLength(args ...interface{}) int {
	return getArgsLength(args) // ERROR `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
}

func someCall() {
	var a = []interface{}{1, 2, 3}
	fmt.Println(checkArgsLength(a...) == getArgsLength(a)) // ERROR `pass \[\]any as any to func getArgsLength func\(args \.\.\.interface\{\}\)`
	fmt.Println(checkArgsLength(a...) == getArgsLength(a...))
}
