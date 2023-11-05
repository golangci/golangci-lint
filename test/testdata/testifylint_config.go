//golangcitest:args -Etestifylint
//golangcitest:config_path testdata/configs/testifylint.yml
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

	assert.Equal(t, predicate, true)
	assert.True(t, resultInt == 1)
	assert.Equal(t, len(arr), 0)
	assert.Error(t, err, io.EOF)
	assert.Nil(t, err)
	assert.Equal(t, resultInt, 42)
	assert.Equal(t, resultFloat, 42.42)
	assert.Equal(t, len(arr), 10)

	assert.True(t, predicate)
	assert.Equal(t, resultInt, 1)
	assert.Empty(t, arr)
	assert.ErrorIs(t, err, io.EOF)
	assert.NoError(t, err) // want "require-error: for error assertions use require"
	assert.Equal(t, 42, resultInt)
	assert.NoErrorf(t, err, "boom!")
	assert.InEpsilon(t, 42.42, resultFloat, 0.0001)
	assert.Len(t, arr, 10)

	require.ErrorIs(t, err, io.EOF)
	require.NoError(t, err)

	t.Run("formatted", func(t *testing.T) {
		assert.Equal(t, predicate, true, "message")
		assert.Equal(t, predicate, true, "message %d", 42)
		assert.Equalf(t, predicate, true, "message")
		assert.Equalf(t, predicate, true, "message %d", 42)
	})

	assert.Equal(t, arr, nil)
	assert.Nil(t, arr)

	go func() {
		if assert.Error(t, err) {
			require.ErrorIs(t, err, io.EOF)
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
	s.Assert().True(b)
}
