package testdata

import (
	"fmt"
)

func uselessFunction(val int) int {
	return val ^ 0
}

func StaticCheckExitFailOnInfo() {
	x := uselessFunction(1)
	fmt.Println(x)
}
