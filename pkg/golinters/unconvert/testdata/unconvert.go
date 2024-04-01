//golangcitest:args -Eunconvert
package testdata

import "io"

// Various explicit conversions of untyped constants
// that cannot be removed.
func _() {
	const (
		_ = byte(0)
		_ = int((real)(0i))
		_ = complex64(complex(1, 2))
		_ = (bool)(true || false)

		PtrSize = 4 << (^uintptr(0) >> 63)
		c0      = uintptr(PtrSize)
		c1      = uintptr((8-PtrSize)/4*2860486313 + (PtrSize-4)/4*33054211828000289)
	)

	i := int64(0)
	_ = i
}

// Make sure we distinguish function calls from
// conversion to function type.
func _() {
	type F func(F) int
	var f F

	_ = F(F(nil)) // want "unnecessary conversion"
	_ = f(F(nil))
}

// Make sure we don't remove explicit conversions that
// prevent fusing floating-point operation.
func _() {
	var f1, f2, f3, ftmp float64
	_ = f1 + float64(f2*f3)
	ftmp = float64(f2 * f3)
	_ = f1 + ftmp
	ftmp = f2 * f3
	_ = f1 + float64(ftmp)

	var c1, c2, c3, ctmp complex128
	_ = c1 + complex128(c2*c3)
	ctmp = complex128(c2 * c3)
	_ = c1 + ctmp
	ctmp = c2 * c3
	_ = c1 + complex128(ctmp)
}

// Basic contains conversion errors for builtin data types
func Basic() {
	var vbool bool
	var vbyte byte
	var vcomplex128 complex128
	var vcomplex64 complex64
	var verror error
	var vfloat32 float32
	var vfloat64 float64
	var vint int
	var vint16 int16
	var vint32 int32
	var vint64 int64
	var vint8 int8
	var vrune rune
	var vstring string
	var vuint uint
	var vuint16 uint16
	var vuint32 uint32
	var vuint64 uint64
	var vuint8 uint8
	var vuintptr uintptr

	_ = bool(vbool)       // want "unnecessary conversion"
	_ = byte(vbyte)       // want "unnecessary conversion"
	_ = error(verror)     // want "unnecessary conversion"
	_ = int(vint)         // want "unnecessary conversion"
	_ = int16(vint16)     // want "unnecessary conversion"
	_ = int32(vint32)     // want "unnecessary conversion"
	_ = int64(vint64)     // want "unnecessary conversion"
	_ = int8(vint8)       // want "unnecessary conversion"
	_ = rune(vrune)       // want "unnecessary conversion"
	_ = string(vstring)   // want "unnecessary conversion"
	_ = uint(vuint)       // want "unnecessary conversion"
	_ = uint16(vuint16)   // want "unnecessary conversion"
	_ = uint32(vuint32)   // want "unnecessary conversion"
	_ = uint64(vuint64)   // want "unnecessary conversion"
	_ = uint8(vuint8)     // want "unnecessary conversion"
	_ = uintptr(vuintptr) // want "unnecessary conversion"

	_ = float32(vfloat32)
	_ = float64(vfloat64)
	_ = complex128(vcomplex128)
	_ = complex64(vcomplex64)

	// Pointers
	_ = (*bool)(&vbool)             // want "unnecessary conversion"
	_ = (*byte)(&vbyte)             // want "unnecessary conversion"
	_ = (*complex128)(&vcomplex128) // want "unnecessary conversion"
	_ = (*complex64)(&vcomplex64)   // want "unnecessary conversion"
	_ = (*error)(&verror)           // want "unnecessary conversion"
	_ = (*float32)(&vfloat32)       // want "unnecessary conversion"
	_ = (*float64)(&vfloat64)       // want "unnecessary conversion"
	_ = (*int)(&vint)               // want "unnecessary conversion"
	_ = (*int16)(&vint16)           // want "unnecessary conversion"
	_ = (*int32)(&vint32)           // want "unnecessary conversion"
	_ = (*int64)(&vint64)           // want "unnecessary conversion"
	_ = (*int8)(&vint8)             // want "unnecessary conversion"
	_ = (*rune)(&vrune)             // want "unnecessary conversion"
	_ = (*string)(&vstring)         // want "unnecessary conversion"
	_ = (*uint)(&vuint)             // want "unnecessary conversion"
	_ = (*uint16)(&vuint16)         // want "unnecessary conversion"
	_ = (*uint32)(&vuint32)         // want "unnecessary conversion"
	_ = (*uint64)(&vuint64)         // want "unnecessary conversion"
	_ = (*uint8)(&vuint8)           // want "unnecessary conversion"
	_ = (*uintptr)(&vuintptr)       // want "unnecessary conversion"
}

// Counter is an int64
type Counter int64

// ID is a typed identifier
type ID string

// Metric is a struct
type Metric struct {
	ID      ID
	Counter Counter
}

// Custom contains conversion errors for builtin data types
func Custom() {
	type Local struct{ id ID }

	var counter Counter
	var id ID
	var m Metric
	var local Local
	var x struct{ id ID }

	_ = Counter(counter)     // want "unnecessary conversion"
	_ = ID(id)               // want "unnecessary conversion"
	_ = Metric(m)            // want "unnecessary conversion"
	_ = Local(local)         // want "unnecessary conversion"
	_ = (struct{ id ID })(x) // want "unnecessary conversion"

	// Pointers
	_ = (*Counter)(&counter)   // want "unnecessary conversion"
	_ = (*ID)(&id)             // want "unnecessary conversion"
	_ = (*Metric)(&m)          // want "unnecessary conversion"
	_ = (*Local)(&local)       // want "unnecessary conversion"
	_ = (*struct{ id ID })(&x) // want "unnecessary conversion"
}

// Interfaces contains conversion errors for interfaces
func Interfaces() {
	var writer io.Writer

	_ = (io.Writer)(writer)   // want "unnecessary conversion"
	_ = (*io.Writer)(&writer) // want "unnecessary conversion"
}

// Constructor is a func type
type Constructor func() ID

// Funcs contains conversion errors for func types
func Funcs() {
	type Local func(ID)
	type Recursive func(Recursive)

	var ctor Constructor
	var local Local
	var recursive Recursive

	_ = Constructor(ctor)    // want "unnecessary conversion"
	_ = Local(local)         // want "unnecessary conversion"
	_ = Recursive(recursive) // want "unnecessary conversion"

	_ = (*Constructor)(&ctor)    // want "unnecessary conversion"
	_ = (*Local)(&local)         // want "unnecessary conversion"
	_ = (*Recursive)(&recursive) // want "unnecessary conversion"
}
