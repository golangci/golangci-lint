//golangcitest:args -Etestifyeqproto
package testdata

import (
	"github.com/anduril/golangci-lint/test/testdata/testifysubpkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

type StructB struct {
	Field2 int
	Field  *testifysubpkg.StructA
}

type StructC struct{
	Field2 int
	Field  testifysubpkg.StructA
}

type StructD struct{}

func TestExample(t *testing.T) {
	sa := &StructB{}
	sc := &StructC{}
	sd := &StructD{}
	assert.Equal(t, sa, sd) // ERROR "call to assert.Equal made with structs that contain proto.Message fields"
	assert.Equal(t, sc, sd) // ERROR "call to assert.Equal made with structs that contain proto.Message fields"
}
