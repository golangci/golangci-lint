//golangcitest:args -Egoerr113
package testdata

import "os"

func SimpleEqual(e1, e2 error) bool {
	return e1 == e2 // want `err113: do not compare errors directly "e1 == e2", use "errors.Is\(e1, e2\)" instead`
}

func SimpleNotEqual(e1, e2 error) bool {
	return e1 != e2 // want `err113: do not compare errors directly "e1 != e2", use "!errors.Is\(e1, e2\)" instead`
}

func CheckGoerr13Import(e error) bool {
	f, err := os.Create("f.txt")
	if err != nil {
		return err == e // want `err113: do not compare errors directly "err == e", use "errors.Is\(err, e\)" instead`
	}
	f.Close()
	return false
}
