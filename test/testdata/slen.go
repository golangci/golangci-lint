package main

import "fmt"

// respect variable position in file except name.
func do() {
	a := []int{}
	if len(a) == 0 {
		return
	}
	_ = a[0]
}

func positive() {
	do()

	a1 := []int{}
	if 0 == len(a1) {
		_ = a1[0]
	}

	a2 := []int{1}
	if len(a2) > 0 {
		_ = a2[0]
	}

	a3 := []int{1}
	if len(a3) < 1 {
		fmt.Println("bad")
	}
	_ = a3[0]

	abc := []int{1, 2, 3, 4}
	for i := range abc {
		_ = abc[i]
	}

	xyz := []int{1, 2, 3}
	for i := 0; i < len(xyz); i++ {
		_ = xyz[i]
	}

	s2 := []int{1, 2, 3}
	if len(s2) == 0 {
		fmt.Println("bad")
	}
	_ = s2[0:]
}

func check(a []int) bool {
	return len(a) == 0
}

func negative() {
	// IndexExpr
	a1 := []int{}
	_ = a1[0] // want `slen: check slice a1 length before accessing`

	// SlicceExpr
	a2 := []int{1, 2, 3}
	_ = a2[0:1] // want `slen: check slice a2 length before accessing`

	// IndexExpr, TODO: handle check if function call
	a3 := []int{1, 2, 3}
	if check(a3) {
		fmt.Println("bad")
	}
	_ = a3[0:1] // want `slen: check slice a3 length before accessing`
}

func main() {
	positive()
	negative()
}
