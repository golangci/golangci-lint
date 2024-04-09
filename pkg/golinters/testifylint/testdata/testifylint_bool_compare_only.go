//golangcitest:args -Etestifylint
//golangcitest:config_path testdata/testifylint_bool_compare_only.yml
package testdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Bool bool

func TestTestifylint(t *testing.T) {
	var predicate bool
	assert.Equal(t, predicate, true) // want "bool-compare: use assert\\.True"
	assert.Equal(t, Bool(predicate), false)
}
