//args: -Egoerr113
package testdata

import "os"

func SimpleEqual(e1, e2 error) bool {
	return e1 == e2 // ERROR `err113: do not compare errors directly, use errors.Is\(\) instead: "e1 == e2"`
}

func SimpleNotEqual(e1, e2 error) bool {
	return e1 != e2 // ERROR `err113: do not compare errors directly, use errors.Is\(\) instead: "e1 != e2"`
}

func CheckGoerr13Import(e error) bool {
	f, err := os.Create("f.txt")
	if err != nil {
		return err == e  // ERROR `err113: do not compare errors directly, use errors.Is\(\) instead: "err == e"`
	}
	f.Close()
	return false
}