//golangcitest:args -Enlreturn
package testdata

func cha() {
	ch := make(chan interface{})
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	select {
	case <-ch:
		return

	case <-ch1:
		{
			a := 1
			_ = a
			{
				a := 1
				_ = a
				return // want "return with no blank line before"
			}

			return
		}

		return

	case <-ch2:
		{
			a := 1
			_ = a
			return // want "return with no blank line before"
		}
		return // want "return with no blank line before"
	}
}

func baz() {
	switch 0 {
	case 0:
		a := 1
		_ = a
		fallthrough // want "fallthrough with no blank line before"
	case 1:
		a := 1
		_ = a
		break // want "break with no blank line before"
	case 2:
		break
	}
}

func foo() int {
	v := []int{}
	for range v {
		return 0
	}

	for range v {
		for range v {
			return 0
		}
		return 0 // want "return with no blank line before"
	}

	o := []int{
		0, 1,
	}

	return o[0]
}

func bar() int {
	o := 1
	if o == 1 {
		if o == 0 {
			return 1
		}
		return 0 // want "return with no blank line before"
	}

	return o
}

func main() {
	return
}

func bugNoAssignSmthHandling() string {
	switch 0 {
	case 0:
		o := struct {
			foo string
		}{
			"foo",
		}
		return o.foo // want "return with no blank line before"

	case 1:
		o := struct {
			foo string
		}{
			"foo",
		}

		return o.foo
	}

	return ""
}

func bugNoExprSmthHandling(string) {
	switch 0 {
	case 0:
		bugNoExprSmthHandling(
			"",
		)
		return // want "return with no blank line before"

	case 1:
		bugNoExprSmthHandling(
			"",
		)

		return
	}
}

func bugNoDeferSmthHandling(string) {
	switch 0 {
	case 0:
		defer bugNoDeferSmthHandling(
			"",
		)
		return // want "return with no blank line before"

	case 1:
		defer bugNoDeferSmthHandling(
			"",
		)

		return
	}
}

func bugNoGoSmthHandling(string) {
	switch 0 {
	case 0:
		go bugNoGoSmthHandling(
			"",
		)
		return // want "return with no blank line before"

	case 1:
		go bugNoGoSmthHandling(
			"",
		)

		return
	}
}
