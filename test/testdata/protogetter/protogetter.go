//golangcitest:args -Eprotogetter
package protogetter

import (
	"fmt"

	"protogetter/proto"
)

type Test struct {
	Embedded *proto.Embedded
}

func testInvalid(t *proto.Test) {
	func(...interface{}) {}(t.B, t.D)      // want "avoid direct access to proto field"
	func(...interface{}) {}(t.GetB(), t.D) // want "avoid direct access to proto field"
	func(...interface{}) {}(t.B, t.GetD()) // want "avoid direct access to proto field"

	_ = t.D                                             // want "avoid direct access to proto field"
	_ = t.F                                             // want "avoid direct access to proto field"
	_ = t.I32                                           // want "avoid direct access to proto field"
	_ = t.I64                                           // want "avoid direct access to proto field"
	_ = t.U32                                           // want "avoid direct access to proto field"
	_ = t.U64                                           // want "avoid direct access to proto field"
	_ = t.T                                             // want "avoid direct access to proto field"
	_ = t.B                                             // want "avoid direct access to proto field"
	_ = t.S                                             // want "avoid direct access to proto field"
	_ = t.Embedded                                      // want "avoid direct access to proto field"
	_ = t.Embedded.S                                    // want "avoid direct access to proto field"
	_ = t.GetEmbedded().S                               // want "avoid direct access to proto field"
	_ = t.Embedded.Embedded                             // want "avoid direct access to proto field"
	_ = t.GetEmbedded().Embedded                        // want "avoid direct access to proto field"
	_ = t.Embedded.Embedded.S                           // want "avoid direct access to proto field"
	_ = t.GetEmbedded().GetEmbedded().S                 // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds                             // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds[0]                          // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds[0].S                        // want "avoid direct access to proto field"
	_ = t.GetRepeatedEmbeddeds()[0].S                   // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds[0].Embedded                 // want "avoid direct access to proto field"
	_ = t.GetRepeatedEmbeddeds()[0].Embedded            // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds[0].Embedded.S               // want "avoid direct access to proto field"
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().S     // want "avoid direct access to proto field"
	_ = t.RepeatedEmbeddeds[t.I64].Embedded.S           // want "avoid direct access to proto field"
	_ = t.GetRepeatedEmbeddeds()[t.I64].GetEmbedded().S // want "avoid direct access to proto field"

	var many []*proto.Test
	manyIndex := 42

	_ = many[0].T                   // want "avoid direct access to proto field"
	_ = many[1].Embedded.S          // want "avoid direct access to proto field"
	_ = many[2].GetEmbedded().S     // want "avoid direct access to proto field"
	_ = many[3].Embedded.Embedded.S // want "avoid direct access to proto field"
	_ = many[manyIndex].S           // want "avoid direct access to proto field"

	test := many[0].Embedded.S == "" || t.Embedded.CustomMethod() == nil || t.S == "" || t.Embedded == nil // want "avoid direct access to proto field"
	_ = test

	other := proto.Other{}
	_ = other.MyMethod(nil).S // want "avoid direct access to proto field"

	ems := t.RepeatedEmbeddeds // want "avoid direct access to proto field"
	_ = ems[len(ems)-1].S      // want "avoid direct access to proto field"

	ch := make(chan string)
	ch <- t.S // want "avoid direct access to proto field"

	for _, v := range t.RepeatedEmbeddeds { // want "avoid direct access to proto field"
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
