//golangcitest:args -Etestifylint
package testdata

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestTestifylint(t *testing.T) {
	var (
		predicate   bool
		resultInt   int
		resultFloat float64
		arr         []string
		err         error
	)

	assert.Equal(t, predicate, true)    // want "bool-compare: use assert\\.True"
	assert.True(t, resultInt == 1)      // want "compares: use assert\\.Equal"
	assert.Equal(t, len(arr), 0)        // want "empty: use assert\\.Empty"
	assert.Error(t, err, io.EOF)        // want "error-is-as: invalid usage of assert\\.Error, use assert\\.ErrorIs instead"
	assert.Nil(t, err)                  // want "error-nil: use assert\\.NoError"
	assert.Equal(t, resultInt, 42)      // want "expected-actual: need to reverse actual and expected values"
	assert.Equal(t, resultFloat, 42.42) // want "float-compare: use assert\\.InEpsilon \\(or InDelta\\)"
	assert.Equal(t, len(arr), 10)       // want "len: use assert\\.Len"

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
	})

	assert.Equal(t, arr, nil) // want "nil-compare: use assert\\.Nil"
	assert.Nil(t, arr)

	go func() {
		if assert.Error(t, err) {
			require.ErrorIs(t, err, io.EOF) // want "go-require: require must only be used in the goroutine running the test function"
		}
	}()
}

type SuiteExample struct {
	suite.Suite
}

func TestSuiteExample(t *testing.T) {
	suite.Run(t, new(SuiteExample))
}

func (s *SuiteExample) TestAll() {
	var b bool
	s.Assert().True(b) // want "suite-extra-assert-call: need to simplify the assertion to s\\.True"
}
