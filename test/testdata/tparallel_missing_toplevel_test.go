//golangcitest:args -Etparallel
package testdata

import (
	"testing"
)

func TestTopLevel(t *testing.T) { // want "TestTopLevel should call t.Parallel on the top level as well as its subtests"
	t.Run("", func(t *testing.T) {
		t.Parallel()
	})
}
