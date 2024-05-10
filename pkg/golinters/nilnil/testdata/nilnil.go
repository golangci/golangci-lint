//golangcitest:args -Enilnil
package testdata

import (
	"bytes"
	"go/token"
	"io"
	"net/http"
	"os"
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

func unsafePtr() (unsafe.Pointer, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr() (uintptr, error) {
	return 0, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0b() (uintptr, error) {
	return 0b0, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0x() (uintptr, error) {
	return 0x00, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0o() (uintptr, error) {
	return 0o000, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
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

func anyType() (any, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m1() (map[int]int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m2() (map[int]*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type mapAlias = map[int]*User

func m3() (mapAlias, error) {
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

type MyError interface {
	error
	Code() string
}

func myError() (*User, MyError) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

// Types.

func structPtrTypeExtPkg() (*os.File, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func primitivePtrTypeExtPkg() (*token.Token, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funcTypeExtPkg() (http.HandlerFunc, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifaceTypeExtPkg() (io.Closer, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type closerAlias = io.Closer

func ifaceTypeAliasedExtPkg() (closerAlias, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
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

type checkerAlias = Checker

func ifaceTypeAliased() (checkerAlias, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type (
	IntegerType    int
	PtrIntegerType *IntegerType
)

func ptrIntegerType() (PtrIntegerType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

// Not checked at all.

func withoutArgs()                                {}
func withoutError1() *User                        { return nil }
func withoutError2() (*User, *User)               { return nil, nil }
func withoutError3() (*User, *User, *User)        { return nil, nil, nil }
func withoutError4() (*User, *User, *User, *User) { return nil, nil, nil, nil }

func invalidOrder() (error, *User)               { return nil, nil }
func withError3rd() (*User, bool, error)         { return nil, false, nil }
func withError4th() (*User, *User, *User, error) { return nil, nil, nil, nil }

func slice() ([]int, error) { return nil, nil }

func strNil() (string, error)   { return "nil", nil }
func strEmpty() (string, error) { return "", nil }

// Valid.

func primitivePtrTypeValid() (*int, error) {
	if false {
		return nil, io.EOF
	}
	return new(int), nil
}

func structPtrTypeValid() (*User, error) {
	if false {
		return nil, io.EOF
	}
	return new(User), nil
}

func unsafePtrValid() (unsafe.Pointer, error) {
	if false {
		return nil, io.EOF
	}
	var i int
	return unsafe.Pointer(&i), nil
}

func uintPtrValid() (uintptr, error) {
	if false {
		return 0, io.EOF
	}
	return 0xc82000c290, nil
}

func channelTypeValid() (ChannelType, error) {
	if false {
		return nil, io.EOF
	}
	return make(ChannelType), nil
}

func funcTypeValid() (FuncType, error) {
	if false {
		return nil, io.EOF
	}
	return func(i int) int {
		return 0
	}, nil
}

func ifaceTypeValid() (io.Reader, error) {
	if false {
		return nil, io.EOF
	}
	return new(bytes.Buffer), nil
}

// Unsupported.

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
