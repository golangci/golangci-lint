package golinters

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	cases := []struct {
		in, out string
		good    bool
	}{
		{"f.go:1:2: text", "", true},
		{"f.go:1:2: text: with: colons", "", true},

		{"f.go:1:2:text wo leading space", "f.go:1:2: text wo leading space", true},

		{"f.go:1:2:", "", false},
		{"f.go:1:2: ", "", false},

		{"f.go:1:2", "f.go:1: 2", true},
		{"f.go:1: text no column", "", true},
		{"f.go:1: text no column: but with colon", "", true},
		{"f.go:1:text no column", "f.go:1: text no column", true},

		{"f.go: no line", "", false},
		{"f.go: 1: text", "", false},

		{"f.go:", "", false},
		{"f.go", "", false},
	}

	lint := TypeCheck{}
	for _, c := range cases {
		i := lint.parseError(errors.New(c.in))
		if !c.good {
			assert.Nil(t, i)
			continue
		}

		assert.NotNil(t, i)

		pos := fmt.Sprintf("%s:%d", i.FilePath(), i.Line())
		if i.Pos.Column != 0 {
			pos += fmt.Sprintf(":%d", i.Pos.Column)
		}
		out := fmt.Sprintf("%s: %s", pos, i.Text)
		expOut := c.out
		if expOut == "" {
			expOut = c.in
		}
		assert.Equal(t, expOut, out)

		assert.Equal(t, "typecheck", i.FromLinter)
	}
}
