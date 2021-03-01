//args: -Ewastedassign
package testdata

import (
	"strings"
)

func p(x int) int {
	return x + 1
}

func typeSwitchNoError(val interface{}, times uint) interface{} {
	switch hoge := val.(type) {
	case int:
		return 12
	case string:
		return strings.Repeat(hoge, int(times))
	default:
		return nil
	}
}

func noUseParamsNoError(params string) int {
	a := 12
	println(a)
	return a
}

func manyif(param int) int {
	println(param)
	useOutOfIf := 1212121 // ERROR "wasted assignment"
	ret := 0
	if false {
		useOutOfIf = 200 // ERROR "reassigned, but never used afterwards"
		return 0
	} else if param == 100 {
		useOutOfIf = 100 // ERROR "wasted assignment"
		useOutOfIf = 201
		useOutOfIf = p(useOutOfIf)
		useOutOfIf += 200 // ERROR "wasted assignment"
	} else {
		useOutOfIf = 100
		useOutOfIf += 100
		useOutOfIf = p(useOutOfIf)
		useOutOfIf += 200 // ERROR "wasted assignment"
	}

	if false {
		useOutOfIf = 200 // ERROR "reassigned, but never used afterwards"
		return 0
	} else if param == 200 {
		useOutOfIf = 100 // ERROR "wasted assignment"
		useOutOfIf = 201
		useOutOfIf = p(useOutOfIf)
		useOutOfIf += 200
	} else {
		useOutOfIf = 100
		useOutOfIf += 100
		useOutOfIf = p(useOutOfIf)
		useOutOfIf += 200
	}
	println(useOutOfIf)
	useOutOfIf = 192
	useOutOfIf += 100
	useOutOfIf += 200 // ERROR "reassigned, but never used afterwards"
	return ret
}

func checkLoopTest() int {
	hoge := 12
	noUse := 1111
	println(noUse)

	noUse = 1111 // ERROR "reassigned, but never used afterwards"
	for {
		if hoge == 14 {
			break
		}
		hoge = hoge + 1
	}
	return hoge
}

func infinity() {
	var i int
	var hoge int
	for {
		hoge = 5 // ERROR "reassigned, but never used afterwards"
	}

	println(i)
	println(hoge)
	return
}

func infinity2() {
	var i int
	var hoge int
	for {
		hoge = 5
		break
	}

	println(i)
	println(hoge)
	return
}
