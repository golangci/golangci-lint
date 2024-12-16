//golangcitest:args -Etestifylint
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
	"fmt"
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type Bool bool

func TestTestifylint(t *testing.T) {
	var (
		predicate   bool
		resultInt   int
		resultFloat float64
		arr         []string
		err         error
	)

	assert.Equal(t, predicate, true)        // want "bool-compare: use assert\\.True"
	assert.Equal(t, Bool(predicate), false) // want "bool-compare: use assert\\.False"
	assert.True(t, resultInt == 1)          // want "compares: use assert\\.Equal"
	assert.Equal(t, len(arr), 0)            // want "empty: use assert\\.Empty"
	assert.Error(t, err, io.EOF)            // want "error-is-as: invalid usage of assert\\.Error, use assert\\.ErrorIs instead"
	assert.Nil(t, err)                      // want "error-nil: use assert\\.NoError"
	assert.Equal(t, resultInt, 42)          // want "expected-actual: need to reverse actual and expected values"
	assert.Equal(t, resultFloat, 42.42)     // want "float-compare: use assert\\.InEpsilon \\(or InDelta\\)"
	assert.Equal(t, len(arr), 10)           // want "len: use assert\\.Len"

	assert.True(t, predicate)
	assert.Equal(t, resultInt, 1) // want "expected-actual: need to reverse actual and expected values"
	assert.Empty(t, arr)
	assert.ErrorIs(t, err, io.EOF) // want "require-error: for error assertions use require"
	assert.NoError(t, err)         // want "require-error: for error assertions use require"
	assert.Equal(t, 42, resultInt)
	assert.InEpsilon(t, 42.42, resultFloat, 0.0001)
	assert.Len(t, arr, 10)

	require.ErrorIs(t, err, io.EOF)
	require.NoError(t, err)

	t.Run("formatted", func(t *testing.T) {
		assert.Equal(t, predicate, true, "message")         // want "bool-compare: use assert\\.True"
		assert.Equal(t, predicate, true, "message %d", 42)  // want "bool-compare: use assert\\.True"
		assert.Equalf(t, predicate, true, "message")        // want "bool-compare: use assert\\.Truef"
		assert.Equalf(t, predicate, true, "message %d", 42) // want "bool-compare: use assert\\.Truef"

		assert.Equal(t, 1, 2, fmt.Sprintf("msg"))     // want "formatter: remove unnecessary fmt\\.Sprintf"
		assert.Equalf(t, 1, 2, "msg with arg", "arg") // want "formatter: assert\\.Equalf call has arguments but no formatting directives"
	})

	assert.Equal(t, arr, nil) // want "nil-compare: use assert\\.Nil"
	assert.Nil(t, arr)

	go func() {
		if assert.Error(t, err) {
			require.ErrorIs(t, err, io.EOF) // want "go-require: require must only be used in the goroutine running the test function"
		}
	}()
}
