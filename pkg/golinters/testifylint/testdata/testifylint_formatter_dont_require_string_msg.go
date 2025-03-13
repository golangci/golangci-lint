//golangcitest:args -Etestifylint
//golangcitest:config_path testdata/testifylint_formatter_dont_require_string_msg.yml
package testdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestifylint(t *testing.T) {
	var b bool
	assert.True(t, b, b)
	assert.True(t, b, "msg %v", 1)
	assert.True(t, b, b, 1) // want "formatter: using msgAndArgs with non-string first element \\(msg\\) causes panic"
}
