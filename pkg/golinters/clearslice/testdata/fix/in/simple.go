//golangcitest:args -Eclearslice
//golangcitest:expected_exitcode 0
package simple

import (
	"fmt"
	"slices"
)

type ReferenceAliasTypeA string
type ReferenceAliasTypeB *int
type PrimitiveAliasType float64

type FlatStruct struct {
	Field1 int
	Field2 float64
	Field3 bool
}

type NestedFlatStruct struct {
	ThingOne FlatStruct
	ThingTwo FlatStruct
}

type FlatLookingStructWithSlice struct {
	Flats []FlatStruct
}

type ReferenceHoldingStruct struct {
	Field1 *FlatStruct
}

type ReferenceHoldingStructUnexported struct {
	field1 *FlatStruct
}

type NestedRefStruct struct {
	Refs []ReferenceHoldingStruct
}

type NestedRefStructUnexported struct {
	refs []ReferenceHoldingStruct
}

type SomeThingWithFlatSliceMember struct {
	Flats []FlatStruct
}

func _() {
	// Safe: slice of primitive types (int)
	s := []int{1, 2, 3}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Safe: slice of primitive types (float)
	s := []float64{1.1, 2.2, 3.3}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Safe: slice of struct free of reference types
	s := []FlatStruct{{1, 1.1, true}, {2, 2.2, false}}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Safe: slice of struct free of reference types
	s := []NestedFlatStruct{
		{ThingOne: FlatStruct{1, 1.1, true}, ThingTwo: FlatStruct{2, 2.2, false}},
	}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of struct that has a slice field (even if slice elements do not contain references)
	s := []FlatLookingStructWithSlice{
		{Flats: []FlatStruct{{1, 1.1, true}, {2, 2.2, false}}},
	}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Safe: slice of alias types that are not reference types
	s := []PrimitiveAliasType{1.1, 2.2, 3.3}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of alias types that are reference types
	s := []ReferenceAliasTypeA{"a", "b", "c"}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of alias types that are reference types
	s := []ReferenceAliasTypeB{new(int), new(int)}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of struct having top level exported fields of reference types
	s := []ReferenceHoldingStruct{{&FlatStruct{1, 1.1, true}}, {&FlatStruct{2, 2.2, false}}}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of struct having unexported fields of reference types
	s := []ReferenceHoldingStructUnexported{{&FlatStruct{1, 1.1, true}}, {&FlatStruct{2, 2.2, false}}}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of struct having nested slices of structs with reference types
	s := []NestedRefStruct{{Refs: []ReferenceHoldingStruct{{&FlatStruct{1, 1.1, true}}, {&FlatStruct{2, 2.2, false}}}}}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of struct having nested slices of structs with unexported reference types
	s := []NestedRefStructUnexported{{refs: []ReferenceHoldingStruct{{&FlatStruct{1, 1.1, true}}, {&FlatStruct{2, 2.2, false}}}}}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of strings
	s := []string{"a", "b", "c"}
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Unsafe: slice of pointers to primitive types
	s := make([]*int, 0, 5)
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Safe: the slice is a member of a struct, but its own elements do not contain references
	o := &SomeThingWithFlatSliceMember{
		Flats: []FlatStruct{{1, 1.1, true}},
	}
	o.Flats = o.Flats[:0]
	fmt.Printf("%+v\n", o)
}

func _() {
	// Unsafe: the slice is a member of a struct, and its own elements do contain references
	o := &NestedRefStructUnexported{
		refs: []ReferenceHoldingStruct{{&FlatStruct{1, 1.1, true}}, {&FlatStruct{2, 2.2, false}}},
	}
	o.refs = o.refs[:0]
	fmt.Printf("%+v\n", o)
}

func _() {
	// Safe pattern that may be prevalent in existing code: clear() directly preceding length adjustment
	s := []*int{new(int), new(int)}
	clear(s)
	s = s[:0]
	fmt.Printf("%+v\n", s)
}

func _() {
	// Recommended pattern: use slices.Delete(s, 0, len(s)) to clear elements up to length
	s := []*int{new(int), new(int)}
	s = slices.Delete(s, 0, len(s))
	fmt.Printf("%+v\n", s)
}
