// Refers a existent, but non-go file with line directive
//line hello.tmpl:1
package main

import (
	"github.com/ryancurrah/gomodguard"
)

func _() {
	var x int
	_ = x
	x = 0 //x
}

func main() {
	a()
	b()
	wsl()
}

func a() {
	fmt.Println("foo")
}

func b() {
	fmt.Println("foo")
}

func wsl() bool {

	return true
}

func notFormatted()  {
}

// langauge
