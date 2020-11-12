package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExcludePatterns(t *testing.T) {
	assert.Equal(t, GetExcludePatterns(nil), DefaultExcludePatterns)

	include := make([]string, 2)
	include[0], include[1] = DefaultExcludePatterns[0].ID, DefaultExcludePatterns[1].ID

	exclude := GetExcludePatterns(include)
	assert.Equal(t, len(exclude), len(DefaultExcludePatterns)-len(include))

	for _, p := range exclude {
		// Not in include.
		for _, i := range include {
			if i == p.ID {
				t.Fatalf("%s can't appear inside include.", p.ID)
			}
		}
		// Must in DefaultExcludePatterns.
		var inDefaultExc bool
		for _, i := range DefaultExcludePatterns {
			if i == p {
				inDefaultExc = true
				break
			}
		}
		assert.True(t, inDefaultExc, fmt.Sprintf("%s must appear inside DefaultExcludePatterns.", p.ID))
	}
}
