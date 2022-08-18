//golangcitest:args -Ewsl
//golangcitest:config_path testdata/configs/wsl.yml
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

	err := multiline(
		"spanning",
		"multiple",
	)
	if err != nil {
		panic(err)
	}

	notErr := multiline(
		"spanning",
		"multiple",
	)
	if err != nil { // ERROR "if statements should only be cuddled with assignments used in the if statement itself"
		panic("not from the line above")
	}

	// This is OK since we use a variable from the line above, even if we don't
	// check it with the if.
	xx := notErr
	if err != nil {
		panic(xx)
	}
}

func multiline(s ...string) error {
	return nil
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

func onelineShouldNotError() error { return nil }

func multilineCase() {
	// Multiline cases
	switch {
	case true,
		false:
		fmt.Println("ok")
	case false ||
		true:
		fmt.Println("ok")
	case true,
		false:
		fmt.Println("ok")
	}
}

func sliceExpr() {
	// Index- and slice expressions.
	var aSlice = []int{1, 2, 3}

	start := 2
	if v := aSlice[start]; v == 1 {
		fmt.Println("ok")
	}

	notOk := 1
	if v := aSlice[start]; v == 1 { // ERROR "if statements should only be cuddled with assignments used in the if statement itself"
		fmt.Println("notOk")
		fmt.Println(notOk)
	}

	end := 2
	if len(aSlice[start:end]) > 2 {
		fmt.Println("ok")
	}
}

func indexExpr() {
	var aMap = map[string]struct{}{"key": {}}

	key := "key"
	if k, ok := aMap[key]; ok {
		fmt.Println(k)
	}

	xxx := "xxx"
	if _, ok := aMap[key]; ok { // ERROR "if statements should only be cuddled with assignments used in the if statement itself"
		fmt.Println("not ok")
		fmt.Println(xxx)
	}
}

func allowTrailing(i int) {
	switch i {
	case 1:
		fmt.Println("one")

	case 2:
		fmt.Println("two")
		// Comments OK too!
	case 3:
		fmt.Println("three")
	}
}

// ExampleSomeOutput simulates an example function.
func ExampleSomeOutput() {
	fmt.Println("Hello, world")

	// Output:
	// Hello, world
}

func IncDecStmt() {
	counter := 0
	for range make([]int, 5) {
		counter++
	}

	type t struct {
		counter int
	}

	x := t{5}

	x.counter--
	if x.counter > 0 {
		fmt.Println("not yet 0")
	}
}

func AnonymousBlock() {
	func(a, b int) { // ERROR "block should not start with a whitespace"

		fmt.Println(a + b)
	}(1, 1)
}

func MultilineComment() {
	if true {
		/*
			Ok to start block with
			a
			long
			multiline
			cmoment
		*/
		fmt.Println("true")
	}
}
