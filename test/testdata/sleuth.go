//args: -Esleuth
package testdata

import "fmt"

func Sleuth() {
	data := []int{1, 2, 3, 4, 5}

	t1 := make([]int, 5)
	for _, tt := range data {
		t1 = append(t1, tt) // ERROR "sleuth found you are trying append to a slice with an initial size"
	}

	fmt.Println(t1)

	t2 := make([]int, 5, 5)
	for i, tt := range data {
		t2[i] = tt // OK
	}

	fmt.Println(t2)

	var t3 []int
	for _, tt := range data {
		t3 = append(t3, tt) // OK
	}

	fmt.Println(t3)

	t4 := make([]int, 0, 5)
	for _, tt := range data {
		t4 = append(t4, tt) // OK
	}

	fmt.Println(t4)
}
