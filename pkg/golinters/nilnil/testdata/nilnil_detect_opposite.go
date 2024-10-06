//golangcitest:args -Enilnil
//golangcitest:config_path testdata/nilnil_detect_opposite.yml
package testdata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"unsafe"
)

func primitivePtrTypeOpposite() (*int, error) {
	if false {
		return nil, io.EOF
	}
	return new(int), errors.New("validation failed") // want "return both a non-nil error and a valid value: use separate returns instead"
}

func structPtrTypeOpposite() (*User, error) {
	if false {
		return nil, io.EOF
	}
	return new(User), fmt.Errorf("invalid %v", 42) // want "return both a non-nil error and a valid value: use separate returns instead"
}

func unsafePtrOpposite() (unsafe.Pointer, error) {
	if false {
		return nil, io.EOF
	}
	var i int
	return unsafe.Pointer(&i), io.EOF // want "return both a non-nil error and a valid value: use separate returns instead"
}

func uintPtrOpposite() (uintptr, error) {
	if false {
		return 0, io.EOF
	}
	return 0xc82000c290, wrap(io.EOF) // want "return both a non-nil error and a valid value: use separate returns instead"
}

func channelTypeOpposite() (ChannelType, error) {
	if false {
		return nil, io.EOF
	}
	return make(ChannelType), fmt.Errorf("wrapped: %w", io.EOF) // want "return both a non-nil error and a valid value: use separate returns instead"
}

func funcTypeOpposite() (FuncType, error) {
	if false {
		return nil, io.EOF
	}
	return func(i int) int { // want "return both a non-nil error and a valid value: use separate returns instead"
		return 0
	}, errors.New("no func type, please")
}

func ifaceTypeOpposite() (io.Reader, error) {
	if false {
		return nil, io.EOF
	}
	return new(bytes.Buffer), new(net.AddrError) // want "return both a non-nil error and a valid value: use separate returns instead"
}

type (
	User             struct{}
	StructPtrType    *User
	PrimitivePtrType *int
	ChannelType      chan int
	FuncType         func(int) int
	Checker          interface{ Check() }
)

func wrap(err error) error { return err }

func structPtr() (*int, error) {
	return nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func structPtrValid() (*int, error) {
	return new(int), nil
}
