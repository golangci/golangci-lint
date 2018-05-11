package testdata

import "fmt"

func GofmtNotSimplified() {
	var x []string
	fmt.Print(x[1:len(x)]) // nolint:megacheck // ERROR "File is not gofmt-ed with -s"
}
