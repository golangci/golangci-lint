package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	assert.True(t, Contains([]string{"val1", "val2", "val3"}, "val2"))
	assert.False(t, Contains([]string{"val1", "val2", "val3"}, "val4"))
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 1, IndexOf([]string{"val1", "val2", "val3"}, "val2"))
	assert.Equal(t, -1, IndexOf([]string{"val1", "val2", "val3"}, "val4"))
}
