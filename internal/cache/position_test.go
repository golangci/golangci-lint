package cache

import (
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestRelativePosition(t *testing.T) {
	mod := &packages.Module{Dir: "/my/path"}

	initialPos := token.Position{
		Filename: filepath.FromSlash("/my/path/foo.go"),
		Offset:   1,
		Line:     2,
		Column:   3,
	}

	pos := RelativePosition(mod, initialPos)

	expected := token.Position{
		Filename: "foo.go",
		Offset:   1,
		Line:     2,
		Column:   3,
	}

	assert.Equal(t, expected, pos)
}

func TestAbsolutionPosition(t *testing.T) {
	mod := &packages.Module{Dir: "/my/other/path"}

	initialPos := token.Position{
		Filename: "foo.go",
		Offset:   1,
		Line:     2,
		Column:   3,
	}

	pos := AbsolutionPosition(mod, initialPos)

	expected := token.Position{
		Filename: filepath.FromSlash("/my/other/path/foo.go"),
		Offset:   1,
		Line:     2,
		Column:   3,
	}

	assert.Equal(t, expected, pos)
}
