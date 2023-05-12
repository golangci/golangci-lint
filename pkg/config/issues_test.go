package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExcludePatterns(t *testing.T) {
	assert.Equal(t, GetExcludePatterns(nil), DefaultExcludePatterns)

	include := []string{DefaultExcludePatterns[0].ID, DefaultExcludePatterns[1].ID}

	exclude := GetExcludePatterns(include)
	assert.Len(t, exclude, len(DefaultExcludePatterns)-len(include))

	for _, p := range exclude {
		assert.NotContains(t, include, p.ID)
		assert.Contains(t, DefaultExcludePatterns, p)
	}
}
