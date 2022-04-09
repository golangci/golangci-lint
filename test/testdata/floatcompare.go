// args: -Efloatcompare
package testdata

import "fmt"

func EqualCompareIfFloats() {
	x, y := 400., 500.
	if 300. == 100. { // ERROR `float comparison found "300. == 100."`
		dummy()
	}
	if x == y { // ERROR `float comparison found "x == y"`
		dummy()
	}
	if 300.+200. == 10. { // ERROR `float comparison found "300.+200. == 10."`
		dummy()
	}
	if 300 == 200 {
		dummy()
	}
}

func NotEqualCompareIfFloats() {
	x, y := 400., 500.
	if 300. != 100. { // ERROR `float comparison found "300. != 100."`

		dummy()
	}
	if x != y { // ERROR `float comparison found "x != y"`
		dummy()
	}
}

func EqualCompareIfCustomType() {
	type number float64
	var x, y number = 300., 400.
	if x == y { // ERROR `float comparison found "x == y"`
		dummy()
	}
}

func GreaterLessCompareIfFloats() {
	if 300. >= 100. { // ERROR `float comparison found "300. >= 100."`
		dummy()
	}
	if 300. <= 100. { // ERROR `float comparison found "300. <= 100."`
		dummy()
	}
	if 300. < 100. { // ERROR `float comparison found "300. < 100."`
		dummy()
	}
	if 300. > 100. { // ERROR `float comparison found "300. > 100."`
		dummy()
	}
}

func SwitchStmtWithFloat() {
	switch 300. { // ERROR "float comparison with switch statement"
	case 100.:
	case 200:
	}
}

func EqualCompareSwitchFloats() {
	switch {
	case 100. == 200.: // ERROR `float comparison found "100. == 200."`
	}
}

func dummy() {
	fmt.Println("dummy()")
}
