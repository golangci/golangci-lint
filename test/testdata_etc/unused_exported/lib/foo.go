package lib

import (
	"fmt"
)

func PublicFunc() {
	privateFunc()
}

func privateFunc() {
	fmt.Println("side effect")
}
