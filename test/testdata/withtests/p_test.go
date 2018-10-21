package withtests

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	if true {
		return
	} else {
		fmt.Printf("test")
	}
}
