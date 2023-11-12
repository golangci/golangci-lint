package gofactory

import "github.com/golangci/golangci-lint/test/testdata/gofactory/nested"

type Struct struct{}

var (
	globalStruct    = nested.Struct{}  // want `Use factory for nested.Struct`
	globalStructPtr = &nested.Struct{} // want `Use factory for nested.Struct`
)

func fn() {
	_ = nested.Struct{}  // want `Use factory for nested.Struct`
	_ = &nested.Struct{} // want `Use factory for nested.Struct`

	_ = []nested.Struct{{}, nested.Struct{}}   // want `Use factory for nested.Struct`
	_ = []*nested.Struct{{}, &nested.Struct{}} // want `Use factory for nested.Struct`

	call(nested.Struct{}) // want `Use factory for nested.Struct`

	_ = []Struct{{}, {}}
}

func call(_ nested.Struct) {}
