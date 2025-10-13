//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package mapsloop

import (
	"iter"
	"maps"
)

var _ = maps.Clone[M] // force "maps" import so that each diagnostic doesn't add one

type M map[int]string

// -- src is map --

func useCopy(dst, src map[int]string) {
	// Replace loop by maps.Copy.
	for key, value := range src {
		// A
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
}

func useCopyGeneric[K comparable, V any, M ~map[K]V](dst, src M) {
	// Replace loop by maps.Copy.
	for key, value := range src {
		// A
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
}

func useCopyNotClone(src map[int]string) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace make(...) by maps.Copy.
	dst := make(map[int]string, len(src))
	// A
	for key, value := range src {
		// B
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
		// C
	}

	// A
	dst = map[int]string{}
	// B
	for key, value := range src {
		// C
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useCopyParen(src map[int]string) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace (make)(...) by maps.Clone.
	dst := (make)(map[int]string, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}

	dst = (map[int]string{})
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useCopy_typesDiffer(src M) {
	// Replace loop but not make(...) as maps.Copy(src) would return wrong type M.
	dst := make(map[int]string, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useCopy_typesDiffer2(src map[int]string) {
	// Replace loop but not make(...) as maps.Copy(src) would return wrong type map[int]string.
	dst := make(M, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useClone_typesDiffer3(src map[int]string) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace loop and make(...) as maps.Clone(src) returns map[int]string
	// which is assignable to M.
	var dst M
	dst = make(M, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useClone_typesDiffer4(src map[int]string) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace loop and make(...) as maps.Clone(src) returns map[int]string
	// which is assignable to M.
	var dst M
	dst = make(M, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

func useClone_generic[Map ~map[K]V, K comparable, V any](src Map) {
	// Clone is tempting but wrong when src may be nil; see #71844.

	// Replace loop and make(...) by maps.Clone
	dst := make(Map, len(src))
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Copy"
	}
	println(dst)
}

// -- src is iter.Seq2 --

func useInsert_assignableToSeq2(dst map[int]string, src func(yield func(int, string) bool)) {
	// Replace loop by maps.Insert because src is assignable to iter.Seq2.
	for k, v := range src {
		dst[k] = v // want "Replace m\\[k\\]=v loop with maps.Insert"
	}
}

func useCollect(src iter.Seq2[int, string]) {
	// Replace loop and make(...) by maps.Collect.
	var dst map[int]string
	dst = make(map[int]string) // A
	// B
	for key, value := range src {
		// C
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Collect"
	}
}

func useInsert_typesDifferAssign(src iter.Seq2[int, string]) {
	// Replace loop and make(...): maps.Collect returns an unnamed map type
	// that is assignable to M.
	var dst M
	dst = make(M)
	// A
	for key, value := range src {
		// B
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Collect"
	}
}

func useInsert_typesDifferDeclare(src iter.Seq2[int, string]) {
	// Replace loop but not make(...) as maps.Collect would return an
	// unnamed map type that would change the type of dst.
	dst := make(M)
	for key, value := range src {
		dst[key] = value // want "Replace m\\[k\\]=v loop with maps.Insert"
	}
}

// -- non-matches --

type isomerOfSeq2 func(yield func(int, string) bool)

func nopeInsertRequiresAssignableToSeq2(dst map[int]string, src isomerOfSeq2) {
	for k, v := range src { // nope: src is not assignable to maps.Insert's iter.Seq2 parameter
		dst[k] = v
	}
}

func nopeSingleVarRange(dst map[int]bool, src map[int]string) {
	for key := range src { // nope: must be "for k, v"
		dst[key] = true
	}
}

func nopeBodyNotASingleton(src map[int]string) {
	var dst map[int]string
	for key, value := range src {
		dst[key] = value
		println() // nope: other things in the loop body
	}
}

// Regression test for https://github.com/golang/go/issues/70815#issuecomment-2581999787.
func nopeAssignmentHasIncrementOperator(src map[int]int) {
	dst := make(map[int]int)
	for k, v := range src {
		dst[k] += v
	}
}

func nopeNotAMap(src map[int]string) {
	var dst []string
	for k, v := range src {
		dst[k] = v
	}
}

func nopeNotAMapGeneric[E any, M ~map[int]E, S ~[]E](src M) {
	var dst S
	for k, v := range src {
		dst[k] = v
	}
}

func nopeHasImplicitWidening(src map[string]int) {
	dst := make(map[string]any)
	for k, v := range src {
		dst[k] = v
	}
}
