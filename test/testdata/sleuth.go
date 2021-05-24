//args: -Esleuth
package testdata

func Sleuth() {
	testdata := []int{1, 2, 3, 4, 5}

	t1 := make([]int, 5)
	for _, tt := range testdata {
		t1 = append(t1, tt) // ERROR "sleuth detects illegal"
	}

	t2 := make([]int, 5, 5)
	for i, tt := range testdata {
		t2[i] = tt // OK
	}

	var t3 []int
	for _, tt := range testdata {
		t3 = append(t3, tt) // OK
	}
}
