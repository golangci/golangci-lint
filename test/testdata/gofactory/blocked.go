//golangcitest:args -Egofactory
//golangcitest:config_path configs/go_factory_only_blocked.yml
package gofactory

import (
	"gofactory/blocked"
	"gofactory/nested"
)

var (
	nestedGlobalStruct    = nested.Struct{}
	nestedGlobalStructPtr = &nested.Struct{}

	blockedGlobalStruct    = blocked.Struct{}
	blockedGlobalStructPtr = &blocked.Struct{}
)

func Blocked() {
	_ = nested.Struct{}
	_ = &nested.Struct{}

	_ = blocked.Struct{}  // want `Use factory for nested.Struct`
	_ = &blocked.Struct{} // want `Use factory for nested.Struct`
}
