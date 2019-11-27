package goanalysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestParseError(t *testing.T) {
	cases := []struct {
		in, out string
		good    bool
	}{
		{"f.go:1:2", "", true},
		{"f.go:1", "", true},
		{"f.go", "", false},
		{"f.go: 1", "", false},
	}

	for _, c := range cases {
		i, _ := parseError(packages.Error{
			Pos: c.in,
			Msg: "msg",
		})
		if !c.good {
			assert.Nil(t, i)
			continue
		}

		assert.NotNil(t, i)

		pos := fmt.Sprintf("%s:%d", i.FilePath(), i.Line())
		if i.Pos.Column != 0 {
			pos += fmt.Sprintf(":%d", i.Pos.Column)
		}
		out := pos
		expOut := c.out
		if expOut == "" {
			expOut = c.in
		}
		assert.Equal(t, expOut, out)

		assert.Equal(t, "typecheck", i.FromLinter)
		assert.Equal(t, "msg", i.Text)
	}
}
