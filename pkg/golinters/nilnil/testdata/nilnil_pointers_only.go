//golangcitest:args -Enilnil
//golangcitest:config_path testdata/nilnil_pointers_only.yml
package testdata

import "unsafe"

type User struct{}

func primitivePtr() (*int, error) {
	return nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func structPtr() (*User, error) {
	return nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func unsafePtr() (unsafe.Pointer, error) {
	return nil, nil
}

func uintPtr0o() (uintptr, error) {
	return 0o000, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func chBi() (chan int, error) {
	return nil, nil
}

func fun() (func(), error) {
	return nil, nil
}

func anyType() (any, error) {
	return nil, nil
}

func m1() (map[int]int, error) {
	return nil, nil
}
