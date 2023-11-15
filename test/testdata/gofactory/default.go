//golangcitest:args -Egofactory
package gofactory

import (
	alias_blocked "gofactory/blocked"
	"gofactory/nested"
)

type Struct struct{}

var (
	defaultGlobalStruct    = nested.Struct{}  // want `Use factory for nested.Struct`
	defaultGlobalStructPtr = &nested.Struct{} // want `Use factory for nested.Struct`
)

func Default() {
	_ = nested.Struct{}  // want `Use factory for nested.Struct`
	_ = &nested.Struct{} // want `Use factory for nested.Struct`

	_ = []nested.Struct{{}, nested.Struct{}}   // want `Use factory for nested.Struct`
	_ = []*nested.Struct{{}, &nested.Struct{}} // want `Use factory for nested.Struct`

	call(nested.Struct{}) // want `Use factory for nested.Struct`

	_ = []Struct{{}, {}}
}

func call(_ nested.Struct) {}

func alias() {
	_ = alias_blocked.Struct{} // want `Use factory for blocked.Struct`
}
