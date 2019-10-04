//args: -Ewsl
//config: linters-settings.wsl.tests=1
package testdata

import (
	"context"
	"fmt"
)

func main() {
	var (
		y = 0
	)
	if y < 1 { // ERROR "if statements should only be cuddled with assignments"
		fmt.Println("tight")
	}

	thisIsNotUsedInIf := true
	if 2 > 1 { // ERROR "if statements should only be cuddled with assignments used in the if statement itself"
		return
	}

	one := 1
	two := 2
	three := 3
	if three == 3 { // ERROR "only one cuddle assignment allowed before if statement"
		fmt.Println("too many cuddled assignments", one, two, thisIsNotUsedInIf)
	}

	var a = "a"
	var b = "b" // ERROR "declarations should never be cuddled"

	if true {
		return
	}
	if false { // ERROR "if statements should only be cuddled with assignments"
		return
	}

	for i := range make([]int, 10) {
		fmt.Println(i)
		fmt.Println(i + i)
		continue // ERROR "branch statements should not be cuddled if block has more than two lines"
	}

	assignOne := a
	fmt.Println(assignOne)
	assignTwo := b // ERROR "assignments should only be cuddled with other assignments"
	fmt.Println(assignTwo)

	_, cf1 := context.WithCancel(context.Background())
	_, cf2 := context.WithCancel(context.Background())
	defer cf1() // ERROR "only one cuddle assignment allowed before defer statement"
	defer cf2()
}

func f1() int {
	x := 1
	return x
}

func f2() int {
	x := 1
	y := 3
	return x + y // ERROR "return statements should not be cuddled if block has more than two lines"
}

func f3() int {
	sum := 0
	for _, v := range []int{2, 4, 8} {
		sum += v
	}

	notSum := 0
	for _, v := range []int{1, 2, 4} { // ERROR "ranges should only be cuddled with assignments used in the iteration"
		sum += v
	}

	return sum + notSum
}
