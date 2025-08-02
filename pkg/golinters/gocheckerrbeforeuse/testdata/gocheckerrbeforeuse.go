//golangcitest:args -Egocheckerrbeforeuse
package testdata

func returns2Values() (int, error) {
	return 0, nil
}

func Negative() {
	i, err := returns2Values() // want "error must be checked right after receiving"

	print(i)

	if err != nil {
		return
	}
}

func Positive() {
	i, err := returns2Values()
	if err != nil {
		return
	}

	print(i)
}
