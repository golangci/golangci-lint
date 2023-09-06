//golangcitest:args -Eprotogolint
package protogolint

import (
	"fmt"

	"protogolint/proto"
)

type Test struct {
	Embedded *proto.Embedded
}

func testInvalid(t *proto.Test) {
	func(...interface{}) {}(t.B, t.D)      // want "proto field read without getter:" "proto field read without getter:"
	func(...interface{}) {}(t.GetB(), t.D) // want "proto field read without getter:"
	func(...interface{}) {}(t.B, t.GetD()) // want "proto field read without getter:"

	_ = t.D                                             // want "proto field read without getter:"
	_ = t.F                                             // want "proto field read without getter:"
	_ = t.I32                                           // want "proto field read without getter:"
	_ = t.I64                                           // want "proto field read without getter:"
	_ = t.U32                                           // want "proto field read without getter:"
	_ = t.U64                                           // want "proto field read without getter:"
	_ = t.T                                             // want "proto field read without getter:"
	_ = t.B                                             // want "proto field read without getter:"
	_ = t.S                                             // want "proto field read without getter:"
	_ = t.Embedded                                      // want "proto field read without getter:"
	_ = t.Embedded.S                                    // want "proto field read without getter:"
	_ = t.GetEmbedded().S                               // want "proto field read without getter:"
	_ = t.Embedded.Embedded                             // want "proto field read without getter:"
	_ = t.GetEmbedded().Embedded                        // want "proto field read without getter:"
	_ = t.Embedded.Embedded.S                           // want "proto field read without getter:"
	_ = t.GetEmbedded().GetEmbedded().S                 // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds                             // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds[0]                          // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds[0].S                        // want "proto field read without getter:"
	_ = t.GetRepeatedEmbeddeds()[0].S                   // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds[0].Embedded                 // want "proto field read without getter:"
	_ = t.GetRepeatedEmbeddeds()[0].Embedded            // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds[0].Embedded.S               // want "proto field read without getter:"
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().S     // want "proto field read without getter:"
	_ = t.RepeatedEmbeddeds[t.I64].Embedded.S           // want "proto field read without getter:"
	_ = t.GetRepeatedEmbeddeds()[t.I64].GetEmbedded().S // want "proto field read without getter:"

	var many []*proto.Test
	manyIndex := 42

	_ = many[0].T                   // want "proto field read without getter:"
	_ = many[1].Embedded.S          // want "proto field read without getter:"
	_ = many[2].GetEmbedded().S     // want "proto field read without getter:"
	_ = many[3].Embedded.Embedded.S // want "proto field read without getter:"
	_ = many[manyIndex].S           // want "proto field read without getter:"

	test := many[0].Embedded.S == "" || t.Embedded.CustomMethod() == nil || t.S == "" || t.Embedded == nil // want "proto field read without getter:" "proto field read without getter:" "proto field read without getter:" "proto field read without getter:"
	_ = test

	other := proto.Other{}
	_ = other.MyMethod(nil).S // want "proto field read without getter:"

	ems := t.RepeatedEmbeddeds // want "proto field read without getter:"
	_ = ems[len(ems)-1].S      // want "proto field read without getter:"

	ch := make(chan string)
	ch <- t.S // want "proto field read without getter:"

	for _, v := range t.RepeatedEmbeddeds { // want "proto field read without getter:"
		_ = v
	}
}

func testValid(t *proto.Test) {
	func(...interface{}) {}(t.GetB(), t.GetD())

	_, t.T = true, true
	_, t.T, _ = true, true, false
	_, _, t.T = true, true, false
	t.T, _ = true, true
	t.D = 2
	t.I32++
	t.I32 += 2

	fmt.Scanf("Test", &t.S, &t.B, &t.T)

	t.D = 1.0
	t.F = 1.0
	t.I32 = 1
	t.I64 = 1
	t.U32 = 1
	t.U64 = 1
	t.T = true
	t.B = []byte{1}
	t.S = "1"
	t.Embedded = &proto.Embedded{}
	t.Embedded.S = "1"
	t.GetEmbedded().S = "1"
	t.Embedded.Embedded = &proto.Embedded{}
	t.GetEmbedded().Embedded = &proto.Embedded{}
	t.Embedded.Embedded.S = "1"
	t.GetEmbedded().GetEmbedded().S = "1"
	t.RepeatedEmbeddeds = []*proto.Embedded{{S: "1"}}

	_ = t.GetD()
	_ = t.GetF()
	_ = t.GetI32()
	_ = t.GetI64()
	_ = t.GetU32()
	_ = t.GetU64()
	_ = t.GetT()
	_ = t.GetB()
	_ = t.GetS()
	_ = t.GetEmbedded()
	_ = t.GetEmbedded().GetS()
	_ = t.GetEmbedded().GetEmbedded()
	_ = t.GetEmbedded().GetEmbedded().GetS()
	_ = t.GetRepeatedEmbeddeds()
	_ = t.GetRepeatedEmbeddeds()[0]
	_ = t.GetRepeatedEmbeddeds()[0].GetS()
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded()
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().GetS()

	other := proto.Other{}
	other.MyMethod(nil).CustomMethod()
	other.MyMethod(nil).GetS()

	var tt Test
	_ = tt.Embedded.GetS()
	_ = tt.Embedded.GetEmbedded().GetS()

	ems := t.GetRepeatedEmbeddeds()
	_ = ems[len(ems)-1].GetS()

	ch := make(chan string)
	ch <- t.GetS()
}
