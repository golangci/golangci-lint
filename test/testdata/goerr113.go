//args: -Egoerr113
package testdata

func SimpleEqual(e1, e2 error) bool {
	return e1 == e2 // ERROR `err113: do not compare errors directly, use errors.Is() instead: "e1 == e2"`
}

func SimpleNotEqual(e1, e2 error) bool {
	return e1 != e2 // ERROR `err113: do not compare errors directly, use errors.Is() instead: "e1 != e2"`
}
