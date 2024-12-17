//golangcitest:args -Etestifylint
package testdata

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Bool bool

func TestTestifylint(t *testing.T) {
	var (
		predicate   bool
		resultInt   int
		resultFloat float64
		arr         []string
		err         error
	)

	assert.True(t, predicate)              // want "bool-compare: use assert\\.True"
	assert.False(t, bool(Bool(predicate))) // want "bool-compare: use assert\\.False"
	assert.Equal(t, resultInt, 1)          // want "compares: use assert\\.Equal"
	assert.Empty(t, arr)                   // want "empty: use assert\\.Empty"
	assert.ErrorIs(t, err, io.EOF)         // want "error-is-as: invalid usage of assert\\.Error, use assert\\.ErrorIs instead"
	assert.NoError(t, err)                 // want "error-nil: use assert\\.NoError"
	assert.Equal(t, 42, resultInt)         // want "expected-actual: need to reverse actual and expected values"
	assert.Equal(t, resultFloat, 42.42)    // want "float-compare: use assert\\.InEpsilon \\(or InDelta\\)"
	assert.Len(t, arr, 10)                 // want "len: use assert\\.Len"

	assert.True(t, predicate)
	assert.Equal(t, 1, resultInt) // want "expected-actual: need to reverse actual and expected values"
	assert.Empty(t, arr)
	assert.ErrorIs(t, err, io.EOF) // want "require-error: for error assertions use require"
	assert.NoError(t, err)         // want "require-error: for error assertions use require"
	assert.Equal(t, 42, resultInt)
	assert.InEpsilon(t, 42.42, resultFloat, 0.0001)
	assert.Len(t, arr, 10)

	require.ErrorIs(t, err, io.EOF)
	require.NoError(t, err)

	t.Run("formatted", func(t *testing.T) {
		assert.True(t, predicate, "message")         // want "bool-compare: use assert\\.True"
		assert.True(t, predicate, "message %d", 42)  // want "bool-compare: use assert\\.True"
		assert.Truef(t, predicate, "message")        // want "bool-compare: use assert\\.Truef"
		assert.Truef(t, predicate, "message %d", 42) // want "bool-compare: use assert\\.Truef"

		assert.Equal(t, 1, 2, "msg")                  // want "formatter: remove unnecessary fmt\\.Sprintf"
		assert.Equalf(t, 1, 2, "msg with arg", "arg") // want "formatter: assert\\.Equalf call has arguments but no formatting directives"
	})

	assert.Nil(t, arr) // want "nil-compare: use assert\\.Nil"
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
	s.True(b) // want "suite-extra-assert-call: need to simplify the assertion to s\\.True"
}

func (s *SuiteExample) TestOne() {

	s.T().Run("subtest", func(t *testing.T) { // want "suite-subtest-run: use s\\.Run to run subtest"

		s.Equal(1, 2) // want "suite-dont-use-pkg: use s\\.Equal"
		s.Equal(1, 2)
	})

	s.Run("subtest", func() {
		s.Equal(1, 2)
	})

	var b bool
	s.True(b) // want "suite-extra-assert-call: need to simplify the assertion to s\\.True"
}
