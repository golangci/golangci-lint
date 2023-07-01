// Refers a existent, but non-go file with line directive
//line hello.tmpl:1
package main

import (
	"golang.org/x/tools/go/analysis"
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

func c(){
	_ = analysis.Analyzer{}
}

func wsl() bool {

	return true
}

func notFormatted()  {
}

// langauge
