//args: -Egocritic
package p

func gocritic() {
	var xs [2048]byte

	// xs -> &xs
	for _, x := range &xs {
		print(x)
	}
}
