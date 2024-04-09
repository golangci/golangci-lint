//golangcitest:args -Enilnil
package testdata

import (
	"io"
	"unsafe"
)

type User struct{}

func primitivePtr() (*int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func structPtr() (*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func emptyStructPtr() (*struct{}, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func anonymousStructPtr() (*struct{ ID string }, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func chBi() (chan int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func chIn() (chan<- int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func chOut() (<-chan int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func fun() (func(), error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funWithArgsAndResults() (func(a, b, c int) (int, int), error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func iface() (interface{}, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m1() (map[int]int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m2() (map[int]*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type Storage struct{}

func (s *Storage) GetUser() (*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifReturn() (*User, error) {
	var s Storage
	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}
	return new(User), nil
}

func forReturn() (*User, error) {
	for {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}
}

func multipleReturn() (*User, error) {
	var s Storage

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	return new(User), nil
}

func nested() {
	_ = func() (*User, error) {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	_, _ = func() (*User, error) {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}()
}

func deeplyNested() {
	_ = func() {
		_ = func() int {
			_ = func() {
				_ = func() (*User, error) {
					_ = func() {}
					_ = func() int { return 0 }
					return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
				}
			}
			return 0
		}
	}
}

type (
	StructPtrType    *User
	PrimitivePtrType *int
	ChannelType      chan int
	FuncType         func(int) int
	Checker          interface{ Check() }
)

func structPtrType() (StructPtrType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func primitivePtrType() (PrimitivePtrType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func channelType() (ChannelType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funcType() (FuncType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifaceType() (Checker, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func withoutArgs()                                {}
func withoutError1() *User                        { return nil }
func withoutError2() (*User, *User)               { return nil, nil }
func withoutError3() (*User, *User, *User)        { return nil, nil, nil }
func withoutError4() (*User, *User, *User, *User) { return nil, nil, nil, nil }

// Unsupported.

func invalidOrder() (error, *User)               { return nil, nil }
func withError3rd() (*User, bool, error)         { return nil, false, nil }
func withError4th() (*User, *User, *User, error) { return nil, nil, nil, nil }
func unsafePtr() (unsafe.Pointer, error)         { return nil, nil }
func uintPtr() (uintptr, error)                  { return 0, nil }
func slice() ([]int, error)                      { return nil, nil }
func ifaceExtPkg() (io.Closer, error)            { return nil, nil }

func implicitNil1() (*User, error) {
	err := (error)(nil)
	return nil, err
}

func implicitNil2() (*User, error) {
	err := io.EOF
	err = nil
	return nil, err
}

func implicitNil3() (*User, error) {
	return nil, wrap(nil)
}
func wrap(err error) error { return err }
