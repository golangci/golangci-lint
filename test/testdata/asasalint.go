package testdata

import "fmt"

func getArgsLength(args ...any) int {
	return len(args)
}

func checkArgsLength(args ...any) int {
	return getArgsLength(args)
}

func someCall() {
	var a = []any{1, 2, 3}
	fmt.Println(checkArgsLength(a...) == getArgsLength(a))
	fmt.Println(checkArgsLength(a...) == getArgsLength(a...))
}
