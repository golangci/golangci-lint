//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package rangeint

import (
	"os"
	os1 "os"
)

func _(i int, s struct{ i int }, slice []int) {
	for i := 0; i < 10; i++ { // want "for loop can be modernized using range over int"
		println(i)
	}
	for j := int(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int8(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int16(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int32(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int64(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := uint8(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := uint16(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := uint32(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := uint64(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int8(0.); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := int8(.0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
	for j := os.FileMode(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}

	{
		var i int
		for i = 0; i < 10; i++ { // want "for loop can be modernized using range over int"
		}
		// NB: no uses of i after loop.
	}
	for i := 0; i < 10; i++ { // want "for loop can be modernized using range over int"
		// i unused within loop
	}
	for i := 0; i < len(slice); i++ { // want "for loop can be modernized using range over int"
		println(slice[i])
	}
	for i := 0; i < len(""); i++ { // want "for loop can be modernized using range over int"
		// NB: not simplified to range ""
	}

	// nope
	for j := .0; j < 10; j++ { // nope: j is a float type
		println(j)
	}
	for j := float64(0); j < 10; j++ { // nope: j is a float type
		println(j)
	}
	for i := 0; i < 10; { // nope: missing increment
	}
	for i := 0; i < 10; i-- { // nope: negative increment
	}
	for i := 0; ; i++ { // nope: missing comparison
	}
	for i := 0; i <= 10; i++ { // nope: wrong comparison
	}
	for ; i < 10; i++ { // nope: missing init
	}
	for s.i = 0; s.i < 10; s.i++ { // nope: not an ident
	}
	for i := 0; i < 10; i++ { // nope: takes address of i
		println(&i)
	}
	for i := 0; i < 10; i++ { // nope: increments i
		i++
	}
	for i := 0; i < 10; i++ { // nope: assigns i
		i = 8
	}

	// The limit expression must be loop invariant;
	// see https://github.com/golang/go/issues/72917
	for i := 0; i < f(); i++ { // nope
	}
	{
		var s struct{ limit int }
		for i := 0; i < s.limit; i++ { // nope: limit is not a const or local var
		}
	}
	{
		const k = 10
		for i := 0; i < k; i++ { // want "for loop can be modernized using range over int"
		}
	}
	{
		var limit = 10
		for i := 0; i < limit; i++ { // want "for loop can be modernized using range over int"
		}
	}
	{
		var limit = 10
		for i := 0; i < limit; i++ { // nope: limit is address-taken
		}
		print(&limit)
	}
	{
		limit := 10
		limit++
		for i := 0; i < limit; i++ { // nope: limit is assigned other than by its declaration
		}
	}
	for i := 0; i < Global; i++ { // nope: limit is an exported global var; may be updated elsewhere
	}
	for i := 0; i < len(table); i++ { // want "for loop can be modernized using range over int"
	}
	{
		s := []string{}
		for i := 0; i < len(s); i++ { // nope: limit is not loop-invariant
			s = s[1:]
		}
	}
	for i := 0; i < len(slice); i++ { // nope: i is incremented within loop
		i += 1
	}
	for Global = 0; Global < 10; Global++ { // nope: loop index is a global variable.
	}
}

var Global int

var table = []string{"hello", "world"}

func f() int { return 0 }

// Repro for part of #71847: ("for range n is invalid if the loop body contains i++"):
func _(s string) {
	var i int                    // (this is necessary)
	for i = 0; i < len(s); i++ { // nope: loop body increments i
		if true {
			i++ // nope
		}
	}
}

// Repro for #71952: for and range loops have different final values
// on i (n and n-1, respectively) so we can't offer the fix if i is
// used after the loop.
func nopePostconditionDiffers() {
	i := 0
	for i = 0; i < 5; i++ {
		println(i)
	}
	println(i) // must print 5, not 4
}

// Non-integer untyped constants need to be explicitly converted to int.
func issue71847d() {
	const limit = 1e3            // float
	for i := 0; i < limit; i++ { // want "for loop can be modernized using range over int"
	}
	for i := int(0); i < limit; i++ { // want "for loop can be modernized using range over int"
	}
	for i := uint(0); i < limit; i++ { // want "for loop can be modernized using range over int"
	}

	const limit2 = 1 + 0i         // complex
	for i := 0; i < limit2; i++ { // want "for loop can be modernized using range over int"
	}
}

func issue72726() {
	var n, kd int
	for i := 0; i < n; i++ { // want "for loop can be modernized using range over int"
		// nope: j will be invisible once it's refactored to 'for j := range min(n-j, kd+1)'
		for j := 0; j < min(n-j, kd+1); j++ { // nope
			_, _ = i, j
		}
	}

	for i := 0; i < i; i++ { // nope
	}

	var i int
	for i = 0; i < i/2; i++ { // nope
	}

	var arr []int
	for i = 0; i < arr[i]; i++ { // nope
	}
}

func todo() {
	for j := os1.FileMode(0); j < 10; j++ { // want "for loop can be modernized using range over int"
		println(j)
	}
}

type T uint
type TAlias = uint

func Fn(a int) T {
	return T(a)
}

func issue73037() {
	var q T
	for a := T(0); a < q; a++ { // want "for loop can be modernized using range over int"
		println(a)
	}
	for a := Fn(0); a < q; a++ {
		println(a)
	}
	var qa TAlias
	for a := TAlias(0); a < qa; a++ { // want "for loop can be modernized using range over int"
		println(a)
	}
	for a := T(0); a < 10; a++ { // want "for loop can be modernized using range over int"
		for b := T(0); b < 10; b++ { // want "for loop can be modernized using range over int"
			println(a, b)
		}
	}
}

func issue75289() {
	// A use of i within a defer may be textually before the loop but runs
	// after, so it should cause the loop to be rejected as a candidate
	// to avoid it observing a different final value of i.
	{
		var i int
		defer func() { println(i) }()
		for i = 0; i < 10; i++ { // nope: i is accessed after the loop (via defer)
		}
	}

	// A use of i within a defer within the loop is also a dealbreaker.
	{
		var i int
		for i = 0; i < 10; i++ { // nope: i is accessed after the loop (via defer)
			defer func() { println(i) }()
		}
	}

	// This (outer) defer is irrelevant.
	defer func() {
		var i int
		for i = 0; i < 10; i++ { // want "for loop can be modernized using range over int"
		}
	}()
}
