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

	blockedGlobalStruct    = blocked.Struct{}  // want `Use factory for blocked.Struct`
	blockedGlobalStructPtr = &blocked.Struct{} // want `Use factory for blocked.Struct`
)

func Blocked() {
	_ = nested.Struct{}
	_ = &nested.Struct{}

	_ = blocked.Struct{}  // want `Use factory for blocked.Struct`
	_ = &blocked.Struct{} // want `Use factory for blocked.Struct`
}
