//golangcitest:args -Egofactory
package testdata

import (
	"net/http"
	alias_blocked "net/http"
)

type Struct struct{}

var (
	defaultGlobalRequest    = http.Request{}  // want `Use factory for http.Request`
	defaultGlobalRequestPtr = &http.Request{} // want `Use factory for http.Request`
)

func Default() {
	_ = http.Request{}  // want `Use factory for http.Request`
	_ = &http.Request{} // want `Use factory for http.Request`

	_ = []http.Request{{}, http.Request{}}   // want `Use factory for http.Request`
	_ = []*http.Request{{}, &http.Request{}} // want `Use factory for http.Request`

	call(http.Request{}) // want `Use factory for http.Request`

	_ = []Struct{{}, {}}
}

func call(_ http.Request) {}

func alias() {
	_ = alias_blocked.Request{} // want `Use factory for http.Request`
}
