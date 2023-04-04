// go:build ruleguard
package ruleguard

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func RangeExprVal(m dsl.Matcher) {
	m.Match(`for _, $_ := range $x { $*_ }`, `for _, $_ = range $x { $*_ }`).
		Where(m["x"].Addressable && m["x"].Type.Size >= 512).
		Report(`$x copy can be avoided with &$x`).
		At(m["x"]).
		Suggest(`&$x`)
}
