package quicktemplate

import (
	"github.com/valyala/bytebufferpool"
)

// ByteBuffer implements io.Writer on top of byte slice.
//
// Recycle byte buffers via AcquireByteBuffer and ReleaseByteBuffer
// in order to reduce memory allocations.
//
// Deprecated: use github.com/valyala/bytebufferpool instead.
type ByteBuffer bytebufferpool.ByteBuffer

// Write implements io.Writer.
func (b *ByteBuffer) Write(p []byte) (int, error) {
	return bb(b).Write(p)
}

// Reset resets the byte buffer.
func (b *ByteBuffer) Reset() {
	bb(b).Reset()
}

// AcquireByteBuffer returns new ByteBuffer from the pool.
//
// Return unneeded buffers to the pool by calling ReleaseByteBuffer
// in order to reduce memory allocations.
func AcquireByteBuffer() *ByteBuffer {
	return (*ByteBuffer)(byteBufferPool.Get())
}

// ReleaseByteBuffer retruns byte buffer to the pool.
//
// Do not access byte buffer after returning it to the pool,
// otherwise data races may occur.
func ReleaseByteBuffer(b *ByteBuffer) {
	byteBufferPool.Put(bb(b))
}

func bb(b *ByteBuffer) *bytebufferpool.ByteBuffer {
	return (*bytebufferpool.ByteBuffer)(b)
}

var byteBufferPool bytebufferpool.Pool
