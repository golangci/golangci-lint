//golangcitest:args -Eprotogetter
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"unsafe"

	"protogetter/proto"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type Test struct {
	Embedded *proto.Embedded
}

func _(t *proto.Test) {
	func(...interface{}) {}(t.B, t.D)      // want `avoid direct access to proto field t\.B, use t\.GetB\(\) instead`
	func(...interface{}) {}(t.GetB(), t.D) // want `avoid direct access to proto field t\.D, use t\.GetD\(\) instead`
	func(...interface{}) {}(t.B, t.GetD()) // want `avoid direct access to proto field t\.B, use t\.GetB\(\) instead`

	_ = t.D                                             // want `avoid direct access to proto field t\.D, use t\.GetD\(\) instead`
	_ = t.F                                             // want `avoid direct access to proto field t\.F, use t\.GetF\(\) instead`
	_ = t.I32                                           // want `avoid direct access to proto field t\.I32, use t\.GetI32\(\) instead`
	_ = t.I64                                           // want `avoid direct access to proto field t\.I64, use t\.GetI64\(\) instead`
	_ = t.U32                                           // want `avoid direct access to proto field t\.U32, use t\.GetU32\(\) instead`
	_ = t.U64                                           // want `avoid direct access to proto field t\.U64, use t\.GetU64\(\) instead`
	_ = t.T                                             // want `avoid direct access to proto field t\.T, use t\.GetT\(\) instead`
	_ = t.B                                             // want `avoid direct access to proto field t\.B, use t\.GetB\(\) instead`
	_ = t.S                                             // want `avoid direct access to proto field t\.S, use t\.GetS\(\) instead`
	_ = t.Embedded                                      // want `avoid direct access to proto field t\.Embedded, use t\.GetEmbedded\(\) instead`
	_ = t.Embedded.S                                    // want `avoid direct access to proto field t\.Embedded\.S, use t\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = t.GetEmbedded().S                               // want `avoid direct access to proto field t\.GetEmbedded\(\)\.S, use t\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = t.Embedded.Embedded                             // want `avoid direct access to proto field t\.Embedded\.Embedded, use t\.GetEmbedded\(\)\.GetEmbedded\(\) instead`
	_ = t.GetEmbedded().Embedded                        // want `avoid direct access to proto field t\.GetEmbedded\(\)\.Embedded, use t\.GetEmbedded\(\)\.GetEmbedded\(\) instead`
	_ = t.Embedded.Embedded.S                           // want `avoid direct access to proto field t\.Embedded\.Embedded\.S, use t\.GetEmbedded\(\)\.GetEmbedded\(\).GetS\(\) instead`
	_ = t.GetEmbedded().GetEmbedded().S                 // want `avoid direct access to proto field t\.GetEmbedded\(\)\.GetEmbedded\(\)\.S, use t\.GetEmbedded\(\)\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = t.RepeatedEmbeddeds                             // want `avoid direct access to proto field t\.RepeatedEmbeddeds, use t\.GetRepeatedEmbeddeds\(\) instead`
	_ = t.RepeatedEmbeddeds[0]                          // want `avoid direct access to proto field t\.RepeatedEmbeddeds, use t\.GetRepeatedEmbeddeds\(\) instead`
	_ = t.RepeatedEmbeddeds[0].S                        // want `avoid direct access to proto field t\.RepeatedEmbeddeds\[0\]\.S, use t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetS\(\) instead`
	_ = t.GetRepeatedEmbeddeds()[0].S                   // want `avoid direct access to proto field t\.GetRepeatedEmbeddeds\(\)\[0\]\.S, use t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetS\(\) instead`
	_ = t.RepeatedEmbeddeds[0].Embedded                 // want `avoid direct access to proto field t\.RepeatedEmbeddeds\[0\]\.Embedded, use t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetEmbedded\(\) instead`
	_ = t.GetRepeatedEmbeddeds()[0].Embedded            // want `avoid direct access to proto field t\.GetRepeatedEmbeddeds\(\)\[0\]\.Embedded, use t\.GetRepeatedEmbeddeds\(\)\[0\]\.GetEmbedded\(\) instead`
	_ = t.RepeatedEmbeddeds[0].Embedded.S               // want `avoid direct access to proto field t\.RepeatedEmbeddeds\[0\]\.Embedded\.S, use t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).GetS\(\) instead`
	_ = t.GetRepeatedEmbeddeds()[0].GetEmbedded().S     // want `avoid direct access to proto field t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).S, use t\.GetRepeatedEmbeddeds\(\)\[0\].GetEmbedded\(\).GetS\(\) instead`
	_ = t.RepeatedEmbeddeds[t.I64].Embedded.S           // want `avoid direct access to proto field t\.RepeatedEmbeddeds\[t.I64\]\.Embedded\.S, use t\.GetRepeatedEmbeddeds\(\)\[t\.GetI64\(\)\].GetEmbedded\(\).GetS\(\) instead`
	_ = t.GetRepeatedEmbeddeds()[t.I64].GetEmbedded().S // want `avoid direct access to proto field t\.GetRepeatedEmbeddeds\(\)\[t\.I64\]\.GetEmbedded\(\)\.S, use t\.GetRepeatedEmbeddeds\(\)\[t\.GetI64\(\)\]\.GetEmbedded\(\).GetS\(\) instead`

	var many []*proto.Test
	manyIndex := 42

	_ = many[0].T                   // want `avoid direct access to proto field many\[0\]\.T, use many\[0\]\.GetT\(\) instead`
	_ = many[1].Embedded.S          // want `avoid direct access to proto field many\[1\]\.Embedded\.S, use many\[1\]\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = many[2].GetEmbedded().S     // want `avoid direct access to proto field many\[2\]\.GetEmbedded\(\)\.S, use many\[2\].GetEmbedded\(\)\.GetS\(\) instead`
	_ = many[3].Embedded.Embedded.S // want `avoid direct access to proto field many\[3\]\.Embedded\.Embedded\.S, use many\[3\].GetEmbedded\(\)\.GetEmbedded\(\)\.GetS\(\) instead`
	_ = many[manyIndex].S           // want `avoid direct access to proto field many\[manyIndex\]\.S, use many\[manyIndex\]\.GetS\(\) instead`

	test := many[0].Embedded.S == "" || t.Embedded.CustomMethod() == nil || t.S == "" || t.Embedded == nil // want `avoid direct access to proto field many\[0\]\.Embedded\.S, use many\[0\]\.GetEmbedded\(\).GetS\(\) instead`
	_ = test

	other := proto.Other{}
	_ = other.MyMethod(nil).S // want `avoid direct access to proto field other\.MyMethod\(nil\)\.S, use other\.MyMethod\(nil\)\.GetS\(\) instead`

	ems := t.RepeatedEmbeddeds // want `avoid direct access to proto field t\.RepeatedEmbeddeds, use t\.GetRepeatedEmbeddeds\(\) instead`
	_ = ems[len(ems)-1].S      // want `avoid direct access to proto field ems\[len\(ems\)-1\]\.S, use ems\[len\(ems\)-1\]\.GetS\(\) instead`

	ch := make(chan string)
	ch <- t.S // want `avoid direct access to proto field t\.S, use t\.GetS\(\) instead`

	for _, v := range t.RepeatedEmbeddeds { // want `avoid direct access to proto field t\.RepeatedEmbeddeds, use t\.GetRepeatedEmbeddeds\(\) instead`
		_ = v
	}

	fn := func(...interface{}) bool { return false }
	fn((*proto.Test)(nil).S) // want `avoid direct access to proto field \(\*proto\.Test\)\(nil\)\.S, use \(\*proto\.Test\)\(nil\)\.GetS\(\) instead`

	var ptrs *[]proto.Test
	_ = (*ptrs)[42].RepeatedEmbeddeds    // want `avoid direct access to proto field \(\*ptrs\)\[42\]\.RepeatedEmbeddeds, use \(\*ptrs\)\[42\].GetRepeatedEmbeddeds\(\) instead`
	_ = (*ptrs)[t.I64].RepeatedEmbeddeds // want `avoid direct access to proto field \(\*ptrs\)\[t\.I64\]\.RepeatedEmbeddeds, use \(\*ptrs\)\[t\.GetI64\(\)\].GetRepeatedEmbeddeds\(\) instead`

	var anyType interface{}
	_ = anyType.(*proto.Test).S // want `avoid direct access to proto field anyType\.\(\*proto\.Test\)\.S, use anyType\.\(\*proto\.Test\)\.GetS\(\) instead`

	t.Embedded.SetS("test")                              // want `avoid direct access to proto field t\.Embedded\.SetS\("test"\), use t\.GetEmbedded\(\)\.SetS\("test"\) instead`
	t.Embedded.SetMap(map[string]string{"test": "test"}) // want `avoid direct access to proto field t\.Embedded\.SetMap\(map\[string\]string{"test": "test"}\), use t\.GetEmbedded\(\)\.SetMap\(map\[string\]string{"test": "test"}\) instead`
}
