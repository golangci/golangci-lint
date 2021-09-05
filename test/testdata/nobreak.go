//args: -Enobreak
package testdata

func infiniteLoop() {
	for { // ERROR "this `for statement` occurs infinite loop"
		println("infinite loop")
	}
}

func noCondInfiniteLoop() {
	for i := 0; ; i++ { // ERROR "this `for statement` occurs infinite loop"
		_ = i
	}
}

func validLoop() {
	for { // OK
		break
	}

	for i := 0; i < 10; i++ { // OK
		_ = i
	}

	for i := 0; ; i++ { // OK
		break
	}
}

func nestedLoop() {
	for { // ERROR "this `for statement` occurs infinite loop"
		for { // OK
			break
		}
	}
}

func nestedLoop2() {
	for { // OK
		for { // OK
			break
		}
		break
	}
}

func nestedLoop3() {
	for { // OK
		for { // ERROR "this `for statement` occurs infinite loop"
			println("infinite loop")
		}
		break
	}
}

func labelLoop() {
loop:
	for { // OK
		break loop
	}
}
