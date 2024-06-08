//golangcitest:args -Etestifylint
//golangcitest:config_path testdata/testifylint_formatter_only.yml
package testdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestifylint(t *testing.T) {
	var err error
	var args []any
	assert.Error(t, err, "Parse(%v) should fail.", args) // want "formatter: use assert\\.Errorf$"

	assert.Equal(t, 1, 2, fmt.Sprintf("msg")) // want "formatter: remove unnecessary fmt\\.Sprintf and use assert\\.Equalf"
	assert.DirExistsf(t, "", "msg with arg", 42)
}
