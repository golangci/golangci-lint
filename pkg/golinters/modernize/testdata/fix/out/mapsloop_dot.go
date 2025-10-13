//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package mapsloop

import . "maps"

var _ = Clone[M] // force "maps" import so that each diagnostic doesn't add one

type M map[int]string

func useCopyDot(dst, src map[int]string) {
	// Replace loop by maps.Copy.
	// want "Replace m\\[k\\]=v loop with maps.Copy"
	Copy(dst, src)
}

func useCloneDot(src map[int]string) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace make(...) by maps.Copy.
	dst := make(map[int]string, len(src))
	// want "Replace m\\[k\\]=v loop with maps.Copy"
	Copy(dst, src)
	println(dst)
}
