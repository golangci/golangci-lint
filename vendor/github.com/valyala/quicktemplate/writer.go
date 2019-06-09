package quicktemplate

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Writer implements auxiliary writer used by quicktemplate functions.
//
// Use AcquireWriter for creating new writers.
type Writer struct {
	e QWriter
	n QWriter
}

// W returns the underlying writer passed to AcquireWriter.
func (qw *Writer) W() io.Writer {
	return qw.n.w
}

// E returns QWriter with enabled html escaping.
func (qw *Writer) E() *QWriter {
	return &qw.e
}

// N returns QWriter without html escaping.
func (qw *Writer) N() *QWriter {
	return &qw.n
}

// AcquireWriter returns new writer from the pool.
//
// Return unneeded writer to the pool by calling ReleaseWriter
// in order to reduce memory allocations.
func AcquireWriter(w io.Writer) *Writer {
	v := writerPool.Get()
	if v == nil {
		qw := &Writer{}
		qw.e.w = &htmlEscapeWriter{}
		v = qw
	}
	qw := v.(*Writer)
	qw.e.w.(*htmlEscapeWriter).w = w
	qw.n.w = w
	return qw
}

// ReleaseWriter returns the writer to the pool.
//
// Do not access released writer, otherwise data races may occur.
func ReleaseWriter(qw *Writer) {
	hw := qw.e.w.(*htmlEscapeWriter)
	hw.w = nil
	qw.e.Reset()
	qw.e.w = hw

	qw.n.Reset()

	writerPool.Put(qw)
}

var writerPool sync.Pool

// QWriter is auxiliary writer used by Writer.
type QWriter struct {
	w   io.Writer
	err error
	b   []byte
}

// Write implements io.Writer.
func (w *QWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.w.Write(p)
	if err != nil {
		w.err = err
	}
	return n, err
}

// Reset resets QWriter to the original state.
func (w *QWriter) Reset() {
	w.w = nil
	w.err = nil
}

// S writes s to w.
func (w *QWriter) S(s string) {
	w.Write(unsafeStrToBytes(s))
}

// Z writes z to w.
func (w *QWriter) Z(z []byte) {
	w.Write(z)
}

// SZ is a synonym to Z.
func (w *QWriter) SZ(z []byte) {
	w.Write(z)
}

// D writes n to w.
func (w *QWriter) D(n int) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendInt(bb.B, int64(n), 10)
	} else {
		w.b = strconv.AppendInt(w.b[:0], int64(n), 10)
		w.Write(w.b)
	}
}

// F writes f to w.
func (w *QWriter) F(f float64) {
	n := int(f)
	if float64(n) == f {
		// Fast path - just int.
		w.D(n)
		return
	}

	// Slow path.
	w.FPrec(f, -1)
}

// FPrec writes f to w using the given floating point precision.
func (w *QWriter) FPrec(f float64, prec int) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = strconv.AppendFloat(bb.B, f, 'f', prec, 64)
	} else {
		w.b = strconv.AppendFloat(w.b[:0], f, 'f', prec, 64)
		w.Write(w.b)
	}
}

// Q writes quoted json-safe s to w.
func (w *QWriter) Q(s string) {
	w.Write(strQuote)
	writeJSONString(w, s)
	w.Write(strQuote)
}

var strQuote = []byte(`"`)

// QZ writes quoted json-safe z to w.
func (w *QWriter) QZ(z []byte) {
	w.Q(unsafeBytesToStr(z))
}

// J writes json-safe s to w.
//
// Unlike Q it doesn't qoute resulting s.
func (w *QWriter) J(s string) {
	writeJSONString(w, s)
}

// JZ writes json-safe z to w.
//
// Unlike Q it doesn't qoute resulting z.
func (w *QWriter) JZ(z []byte) {
	w.J(unsafeBytesToStr(z))
}

// V writes v to w.
func (w *QWriter) V(v interface{}) {
	fmt.Fprintf(w, "%v", v)
}

// U writes url-encoded s to w.
func (w *QWriter) U(s string) {
	bb, ok := w.w.(*ByteBuffer)
	if ok {
		bb.B = appendURLEncode(bb.B, s)
	} else {
		w.b = appendURLEncode(w.b[:0], s)
		w.Write(w.b)
	}
}

// UZ writes url-encoded z to w.
func (w *QWriter) UZ(z []byte) {
	w.U(unsafeBytesToStr(z))
}
