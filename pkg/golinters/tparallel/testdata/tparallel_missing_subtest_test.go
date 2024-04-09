//golangcitest:args -Etparallel
package testdata

import (
	"testing"
)

func TestSubtests(t *testing.T) { // want "TestSubtests's subtests should call t.Parallel"
	t.Parallel()

	t.Run("", func(t *testing.T) {
	})
}
