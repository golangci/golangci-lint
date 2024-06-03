//golangcitest:args -Ewhyvarscope
package testdata

func zero() int {
	return 0
}

func main() {
	if z := zero(); z == 0 { // want "variable z can be removed and use assignee directly"
		println("z is 0")
	}

	if z := zero(); z == 0 {
		println(z)
	}
}
