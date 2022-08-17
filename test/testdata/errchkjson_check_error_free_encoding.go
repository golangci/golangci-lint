//golangcitest:args -Eerrchkjson
//golangcitest:config_path testdata/configs/errchkjson_check_error_free_encoding.yml
package testdata

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"unsafe"
)

type marshalText struct{}

func (mt marshalText) MarshalText() ([]byte, error) {
	return []byte(`mt`), nil
}

var _ encoding.TextMarshaler = marshalText(struct{}{})

// JSONMarshalSafeTypes contains a multitude of test cases to marshal different combinations of types to JSON,
// that are safe, that is, they will never return an error, if these types are marshaled to JSON.
func JSONMarshalSafeTypes() {
	var err error

	_, _ = json.Marshal(nil)   // nil is safe
	json.Marshal(nil)          // nil is safe
	_, err = json.Marshal(nil) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	_, _ = json.MarshalIndent(nil, "", " ")   // nil is safe
	json.MarshalIndent(nil, "", " ")          // nil is safe
	_, err = json.MarshalIndent(nil, "", " ") // ERROR "Error return value of `encoding/json.MarshalIndent` is checked but passed argument is safe"
	_ = err

	enc := json.NewEncoder(ioutil.Discard)
	_ = enc.Encode(nil) // ERROR "Error return value of `\\([*]encoding/json.Encoder\\).Encode` is not checked"
	enc.Encode(nil)     // ERROR "Error return value of `\\([*]encoding/json.Encoder\\).Encode` is not checked"
	err = enc.Encode(nil)
	_ = err

	var b bool
	_, _ = json.Marshal(b)   // bool is safe
	_, err = json.Marshal(b) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var i int
	_, _ = json.Marshal(i)   // int is safe
	_, err = json.Marshal(i) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var i8 int8
	_, _ = json.Marshal(i8)   // int8 is safe
	_, err = json.Marshal(i8) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var i16 int16
	_, _ = json.Marshal(i16)   // int16 is safe
	_, err = json.Marshal(i16) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var i32 int32
	_, _ = json.Marshal(i32)   // int32 / rune is safe
	_, err = json.Marshal(i32) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var i64 int64
	_, _ = json.Marshal(i64)   // int64 is safe
	_, err = json.Marshal(i64) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ui uint
	_, _ = json.Marshal(ui)   // uint is safe
	_, err = json.Marshal(ui) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ui8 uint8
	_, _ = json.Marshal(ui8)   // uint8 / byte is safe
	_, err = json.Marshal(ui8) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ui16 uint16
	_, _ = json.Marshal(ui16)   // uint16 is safe
	_, err = json.Marshal(ui16) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ui32 uint32
	_, _ = json.Marshal(ui32)   // uint32 / rune is safe
	_, err = json.Marshal(ui32) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ui64 uint64
	_, _ = json.Marshal(ui64)   // uint64 is safe
	_, err = json.Marshal(ui64) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var uiptr uintptr
	_, _ = json.Marshal(uiptr)   // uintptr is safe
	_, err = json.Marshal(uiptr) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var str string
	_, _ = json.Marshal(str)   // string is safe
	_, err = json.Marshal(str) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var strSlice []string
	_, _ = json.Marshal(strSlice)   // []string is safe
	_, err = json.Marshal(strSlice) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var intSlice []int
	_, _ = json.Marshal(intSlice)   // []int is safe
	_, err = json.Marshal(intSlice) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var boolSlice []bool
	_, _ = json.Marshal(boolSlice)   // []bool is safe
	_, err = json.Marshal(boolSlice) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var strArray [10]string
	_, _ = json.Marshal(strArray)   // [10]string is safe
	_, err = json.Marshal(strArray) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var intArray [10]int
	_, _ = json.Marshal(intArray)   // [10]int is safe
	_, err = json.Marshal(intArray) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var boolArray [10]bool
	_, _ = json.Marshal(boolArray)   // [10]bool is safe
	_, err = json.Marshal(boolArray) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var basicStruct struct {
		Bool    bool
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32 // also rune
		Int64   int64
		Uint    uint
		Uint8   uint8 // also byte
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Uintptr uintptr
		String  string
	}
	_, _ = json.Marshal(basicStruct)   // struct containing only safe basic types is safe
	_, err = json.Marshal(basicStruct) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var ptrStruct struct {
		Bool    *bool
		Int     *int
		Int8    *int8
		Int16   *int16
		Int32   *int32
		Int64   *int64
		Uint    *uint
		Uint8   *uint8
		Uint16  *uint16
		Uint32  *uint32
		Uint64  *uint64
		Uintptr *uintptr
		String  *string
	}
	_, _ = json.Marshal(ptrStruct)   // struct containing pointer to only safe basic types is safe
	_, err = json.Marshal(ptrStruct) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapStrStr map[string]string
	_, _ = json.Marshal(mapStrStr)   // map[string]string is safe
	_, err = json.Marshal(mapStrStr) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapStrInt map[string]int
	_, _ = json.Marshal(mapStrInt)   // map[string]int is safe
	_, err = json.Marshal(mapStrInt) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapStrBool map[string]bool
	_, _ = json.Marshal(mapStrBool)   // map[string]bool is safe
	_, err = json.Marshal(mapStrBool) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapIntStr map[int]string
	_, _ = json.Marshal(mapIntStr)   // map[int]string is safe
	_, err = json.Marshal(mapIntStr) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapIntInt map[int]int
	_, _ = json.Marshal(mapIntInt)   // map[int]int is safe
	_, err = json.Marshal(mapIntInt) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	var mapIntBool map[int]bool
	_, _ = json.Marshal(mapIntBool)   // map[int]bool is safe
	_, err = json.Marshal(mapIntBool) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err

	type innerStruct struct {
		Bool   bool
		Int    int
		String string

		StrSlice  []string
		IntSlice  []int
		BoolSlice []bool

		StrArray  [10]string
		IntArray  [10]int
		BoolArray [10]bool

		MapStrStr  map[string]string
		MapStrInt  map[string]int
		MapStrBool map[string]bool

		MapIntStr  map[int]string
		MapIntInt  map[int]int
		MapIntBool map[int]bool
	}
	var outerStruct struct {
		Bool   bool
		Int    int
		String string

		StrSlice  []string
		IntSlice  []int
		BoolSlice []bool

		StrArray  [10]string
		IntArray  [10]int
		BoolArray [10]bool

		MapStrStr  map[string]string
		MapStrInt  map[string]int
		MapStrBool map[string]bool

		MapIntStr  map[int]string
		MapIntInt  map[int]int
		MapIntBool map[int]bool

		InnerStruct innerStruct
	}
	_, _ = json.Marshal(outerStruct)   // struct with only safe types
	_, err = json.Marshal(outerStruct) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err
}

type (
	structKey                      struct{ id int }
	ExportedUnsafeAndInvalidStruct struct { // unsafe unexported but omitted
		F64                  float64
		F64Ptr               *float64
		F64Slice             []float64
		F64Array             [10]float64
		MapStrF64            map[string]float64
		MapEIStr             map[interface{}]string
		Number               json.Number
		NumberPtr            *json.Number
		NumberSlice          []json.Number
		MapNumberStr         map[json.Number]string
		Ei                   interface{}
		Stringer             fmt.Stringer
		Mt                   marshalText
		MapMarshalTextString map[marshalText]string

		C128         complex128
		C128Ptr      *complex128
		C128Slice    []complex128
		C128Array    [10]complex128
		MapBoolStr   map[bool]string
		MapF64Str    map[float64]string
		F            func()
		Ch           chan struct{}
		UnsafePtr    unsafe.Pointer
		MapStructStr map[structKey]string
	}
)

// JSONMarshalSaveStructWithUnexportedFields contains a struct with unexported, unsafe fields.
func JSONMarshalSaveStructWithUnexportedFields() {
	var err error

	var unexportedInStruct struct {
		Bool bool // safe exported

		f64                  float64                         // unsafe unexported
		f64Ptr               *float64                        // unsafe unexported
		f64Slice             []float64                       // unsafe unexported
		f64Array             [10]float64                     // unsafe unexported
		mapStrF64            map[string]float64              // unsafe unexported
		mapEIStr             map[interface{}]string          // unsafe unexported
		number               json.Number                     // unsafe unexported
		numberPtr            *json.Number                    // unsafe unexported
		numberSlice          []json.Number                   // unsafe unexported
		mapNumberStr         map[json.Number]string          // unsafe unexported
		ei                   interface{}                     // unsafe unexported
		stringer             fmt.Stringer                    // unsafe unexported
		mt                   marshalText                     // unsafe unexported
		mapMarshalTextString map[marshalText]string          // unsafe unexported
		unexportedStruct     ExportedUnsafeAndInvalidStruct  // unsafe unexported
		unexportedStructPtr  *ExportedUnsafeAndInvalidStruct // unsafe unexported

		c128         complex128           // invalid unexported
		c128Slice    []complex128         // invalid unexported
		c128Array    [10]complex128       // invalid unexported
		mapBoolStr   map[bool]string      // invalid unexported
		mapF64Str    map[float64]string   // invalid unexported
		f            func()               // invalid unexported
		ch           chan struct{}        // invalid unexported
		unsafePtr    unsafe.Pointer       // invalid unexported
		mapStructStr map[structKey]string // invalid unexported
	}
	_ = unexportedInStruct.f64
	_ = unexportedInStruct.f64Ptr
	_ = unexportedInStruct.f64Slice
	_ = unexportedInStruct.f64Array
	_ = unexportedInStruct.mapStrF64
	_ = unexportedInStruct.mapEIStr
	_ = unexportedInStruct.number
	_ = unexportedInStruct.numberPtr
	_ = unexportedInStruct.numberSlice
	_ = unexportedInStruct.mapNumberStr
	_ = unexportedInStruct.ei
	_ = unexportedInStruct.stringer
	_ = unexportedInStruct.mt
	_ = unexportedInStruct.mapMarshalTextString
	_ = unexportedInStruct.unexportedStruct
	_ = unexportedInStruct.unexportedStructPtr

	_ = unexportedInStruct.c128
	_ = unexportedInStruct.c128Slice
	_ = unexportedInStruct.c128Array
	_ = unexportedInStruct.mapBoolStr
	_ = unexportedInStruct.mapF64Str
	_ = unexportedInStruct.f
	_ = unexportedInStruct.ch
	_ = unexportedInStruct.unsafePtr
	_ = unexportedInStruct.mapStructStr[structKey{1}]
	_, _ = json.Marshal(unexportedInStruct)   // struct containing unsafe but unexported fields is safe
	_, err = json.Marshal(unexportedInStruct) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err
}

// JSONMarshalSaveStructWithOmittedFields contains a struct with omitted, unsafe fields.
func JSONMarshalSaveStructWithOmittedFields() {
	var err error

	var omitInStruct struct {
		Bool bool // safe exported

		F64                  float64                         `json:"-"` // unsafe exported but omitted
		F64Ptr               *float64                        `json:"-"` // unsafe exported but omitted
		F64Slice             []float64                       `json:"-"` // unsafe exported but omitted
		F64Array             [10]float64                     `json:"-"` // unsafe exported but omitted
		MapStrF64            map[string]float64              `json:"-"` // unsafe exported but omitted
		MapEIStr             map[interface{}]string          `json:"-"` // unsafe exported but omitted
		Number               json.Number                     `json:"-"` // unsafe exported but omitted
		NumberPtr            *json.Number                    `json:"-"` // unsafe exported but omitted
		NumberSlice          []json.Number                   `json:"-"` // unsafe exported but omitted
		MapNumberStr         map[json.Number]string          `json:"-"` // unsafe exported but omitted
		Ei                   interface{}                     `json:"-"` // unsafe exported but omitted
		Stringer             fmt.Stringer                    `json:"-"` // unsafe exported but omitted
		Mt                   marshalText                     `json:"-"` // unsafe exported but omitted
		MapMarshalTextString map[marshalText]string          `json:"-"` // unsafe exported but omitted
		ExportedStruct       ExportedUnsafeAndInvalidStruct  `json:"-"` // unsafe exported but omitted
		ExportedStructPtr    *ExportedUnsafeAndInvalidStruct `json:"-"` // unsafe exported but omitted

		C128         complex128           `json:"-"` // invalid exported but omitted
		C128Slice    []complex128         `json:"-"` // invalid exported but omitted
		C128Array    [10]complex128       `json:"-"` // invalid exported but omitted
		MapBoolStr   map[bool]string      `json:"-"` // invalid exported but omitted
		MapF64Str    map[float64]string   `json:"-"` // invalid exported but omitted
		F            func()               `json:"-"` // invalid exported but omitted
		Ch           chan struct{}        `json:"-"` // invalid exported but omitted
		UnsafePtr    unsafe.Pointer       `json:"-"` // invalid exported but omitted
		MapStructStr map[structKey]string `json:"-"` // invalid exported but omitted
	}
	_ = omitInStruct.MapStructStr[structKey{1}]
	_, _ = json.Marshal(omitInStruct)   // struct containing unsafe but omitted, exported fields is safe
	_, err = json.Marshal(omitInStruct) // ERROR "Error return value of `encoding/json.Marshal` is checked but passed argument is safe"
	_ = err
}

// JSONMarshalUnsafeTypes contains a multitude of test cases to marshal different combinations of types to JSON,
// that can potentially lead to json.Marshal returning an error.
func JSONMarshalUnsafeTypes() {
	var err error

	var f32 float32
	json.Marshal(f32)          // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, _ = json.Marshal(f32)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, err = json.Marshal(f32) // err is checked
	_ = err

	var f64 float64
	_, _ = json.Marshal(f64)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float64` found"
	_, err = json.Marshal(f64) // err is checked
	_ = err

	var f32Slice []float32
	_, _ = json.Marshal(f32Slice)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, err = json.Marshal(f32Slice) // err is checked
	_ = err

	var f64Slice []float64
	_, _ = json.Marshal(f64Slice)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float64` found"
	_, err = json.Marshal(f64Slice) // err is checked
	_ = err

	var f32Array [10]float32
	_, _ = json.Marshal(f32Array)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, err = json.Marshal(f32Array) // err is checked
	_ = err

	var f64Array [10]float64
	_, _ = json.Marshal(f64Array)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float64` found"
	_, err = json.Marshal(f64Array) // err is checked
	_ = err

	var structPtrF32 struct {
		F32 *float32
	}
	_, _ = json.Marshal(structPtrF32)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, err = json.Marshal(structPtrF32) // err is checked
	_ = err

	var structPtrF64 struct {
		F64 *float64
	}
	_, _ = json.Marshal(structPtrF64)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float64` found"
	_, err = json.Marshal(structPtrF64) // err is checked
	_ = err

	var mapStrF32 map[string]float32
	_, _ = json.Marshal(mapStrF32)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float32` found"
	_, err = json.Marshal(mapStrF32) // err is checked
	_ = err

	var mapStrF64 map[string]float64
	_, _ = json.Marshal(mapStrF64)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `float64` found"
	_, err = json.Marshal(mapStrF64) // err is checked
	_ = err

	var mapEIStr map[interface{}]string
	_, _ = json.Marshal(mapEIStr)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `interface{}` as map key found"
	_, err = json.Marshal(mapEIStr) // err is checked
	_ = err

	var mapStringerStr map[fmt.Stringer]string
	_, _ = json.Marshal(mapStringerStr)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `fmt.Stringer` as map key found"
	_, err = json.Marshal(mapStringerStr) // err is checked
	_ = err

	var number json.Number
	_, _ = json.Marshal(number)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `encoding/json.Number` found"
	_, err = json.Marshal(number) // err is checked
	_ = err

	var numberSlice []json.Number
	_, _ = json.Marshal(numberSlice)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `encoding/json.Number` found"
	_, err = json.Marshal(numberSlice) // err is checked
	_ = err

	var mapNumberStr map[json.Number]string
	_, _ = json.Marshal(mapNumberStr)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `encoding/json.Number` as map key found"
	_, err = json.Marshal(mapNumberStr) // err is checked
	_ = err

	var mapStrNumber map[string]json.Number
	_, _ = json.Marshal(mapStrNumber)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `encoding/json.Number` found"
	_, err = json.Marshal(mapStrNumber) // err is checked
	_ = err

	var ei interface{}
	_, _ = json.Marshal(ei)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `interface{}` found"
	_, err = json.Marshal(ei) // err is checked
	_ = err

	var eiptr *interface{}
	_, _ = json.Marshal(eiptr)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `*interface{}` found"
	_, err = json.Marshal(eiptr) // err is checked
	_ = err

	var stringer fmt.Stringer
	_, _ = json.Marshal(stringer)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `fmt.Stringer` found"
	_, err = json.Marshal(stringer) // err is checked
	_ = err

	var structWithEmptyInterface struct {
		EmptyInterface interface{}
	}
	_, _ = json.Marshal(structWithEmptyInterface)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `interface{}` found"
	_, err = json.Marshal(structWithEmptyInterface) // err is checked
	_ = err

	var mt marshalText
	_, _ = json.Marshal(mt)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `[a-z-]+.marshalText` found"
	_, err = json.Marshal(mt) // err is checked
	_ = err

	var mapMarshalTextString map[marshalText]string
	_, _ = json.Marshal(mapMarshalTextString)   // ERROR "Error return value of `encoding/json.Marshal` is not checked: unsafe type `[a-z-]+.marshalText` as map key found"
	_, err = json.Marshal(mapMarshalTextString) // err is checked
	_ = err
}

// JSONMarshalInvalidTypes contains a multitude of test cases to marshal different combinations of types to JSON,
// that are invalid and not supported by json.Marshal, that is they will always return an error, if these types used
// with json.Marshal.
func JSONMarshalInvalidTypes() {
	var err error

	var c64 complex64
	json.Marshal(c64)          // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_, _ = json.Marshal(c64)   // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_, err = json.Marshal(c64) // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_ = err

	var c128 complex128
	_, _ = json.Marshal(c128)   // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_, err = json.Marshal(c128) // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_ = err

	var sliceC64 []complex64
	_, _ = json.Marshal(sliceC64)   // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_, err = json.Marshal(sliceC64) // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_ = err

	var sliceC128 []complex128
	_, _ = json.Marshal(sliceC128)   // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_, err = json.Marshal(sliceC128) // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_ = err

	var arrayC64 []complex64
	_, _ = json.Marshal(arrayC64)   // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_, err = json.Marshal(arrayC64) // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_ = err

	var arrayC128 []complex128
	_, _ = json.Marshal(arrayC128)   // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_, err = json.Marshal(arrayC128) // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_ = err

	var structPtrC64 struct {
		C64 *complex64
	}
	_, _ = json.Marshal(structPtrC64)   // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_, err = json.Marshal(structPtrC64) // ERROR "`encoding/json.Marshal` for unsupported type `complex64` found"
	_ = err

	var structPtrC128 struct {
		C128 *complex128
	}
	_, _ = json.Marshal(structPtrC128)   // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_, err = json.Marshal(structPtrC128) // ERROR "`encoding/json.Marshal` for unsupported type `complex128` found"
	_ = err

	var mapBoolStr map[bool]string
	_, _ = json.Marshal(mapBoolStr)   // ERROR "`encoding/json.Marshal` for unsupported type `bool` as map key found"
	_, err = json.Marshal(mapBoolStr) // ERROR "`encoding/json.Marshal` for unsupported type `bool` as map key found"
	_ = err

	var mapF32Str map[float32]string
	_, _ = json.Marshal(mapF32Str)   // ERROR "`encoding/json.Marshal` for unsupported type `float32` as map key found"
	_, err = json.Marshal(mapF32Str) // ERROR "`encoding/json.Marshal` for unsupported type `float32` as map key found"
	_ = err

	var mapF64Str map[float64]string
	_, _ = json.Marshal(mapF64Str)   // ERROR "`encoding/json.Marshal` for unsupported type `float64` as map key found"
	_, err = json.Marshal(mapF64Str) // ERROR "`encoding/json.Marshal` for unsupported type `float64` as map key found"
	_ = err

	var mapC64Str map[complex64]string
	_, _ = json.Marshal(mapC64Str)   // ERROR "`encoding/json.Marshal` for unsupported type `complex64` as map key found"
	_, err = json.Marshal(mapC64Str) // ERROR "`encoding/json.Marshal` for unsupported type `complex64` as map key found"
	_ = err

	var mapC128Str map[complex128]string
	_, _ = json.Marshal(mapC128Str)   // ERROR "`encoding/json.Marshal` for unsupported type `complex128` as map key found"
	_, err = json.Marshal(mapC128Str) // ERROR "`encoding/json.Marshal` for unsupported type `complex128` as map key found"
	_ = err

	mapStructStr := map[structKey]string{structKey{1}: "str"}
	_, _ = json.Marshal(mapStructStr)   // ERROR "`encoding/json.Marshal` for unsupported type `[a-z-]+.structKey` as map key found"
	_, err = json.Marshal(mapStructStr) // ERROR "`encoding/json.Marshal` for unsupported type `[a-z-]+.structKey` as map key found"
	_ = err

	f := func() {}
	_, _ = json.Marshal(f)   // ERROR "`encoding/json.Marshal` for unsupported type `func\\(\\)` found"
	_, err = json.Marshal(f) // ERROR "`encoding/json.Marshal` for unsupported type `func\\(\\)` found"
	_ = err

	var ch chan struct{} = make(chan struct{})
	_, _ = json.Marshal(ch)   // ERROR "`encoding/json.Marshal` for unsupported type `chan struct{}` found"
	_, err = json.Marshal(ch) // ERROR "`encoding/json.Marshal` for unsupported type `chan struct{}` found"
	_ = err

	var unsafePtr unsafe.Pointer
	_, _ = json.Marshal(unsafePtr)   // ERROR "`encoding/json.Marshal` for unsupported type `unsafe.Pointer` found"
	_, err = json.Marshal(unsafePtr) // ERROR "`encoding/json.Marshal` for unsupported type `unsafe.Pointer` found"
	_ = err
}

// NotJSONMarshal contains other go ast node types, that are not considered by errchkjson
func NotJSONMarshal() {
	s := fmt.Sprintln("I am not considered by errchkjson")
	_ = s
	f := func() bool { return false }
	_ = f()
}
