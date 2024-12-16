//golangcitest:args -Eerr113
//golangcitest:expected_exitcode 0
package testdata

import "os"

func SimpleEqual(e1, e2 error) bool {
	return e1 == e2
}

func SimpleNotEqual(e1, e2 error) bool {
	return e1 != e2
}

func CheckGoerr13Import(e error) bool {
	f, err := os.Create("f.txt")
	if err != nil {
		return err == e
	}
	f.Close()
	return false
}
