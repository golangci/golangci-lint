package testdata

import "fmt"

func gofmtNotSimplified() {
	var x []string
	fmt.Print(x[1:len(x)]) // ERROR "File is not gofmt-ed with -s"
}
