//golangcitest:args -Etparallel
//golangcitest:expected_exitcode 0
package testdata

import (
	"testing"
)

func TestValidHappyPath(t *testing.T) {
	t.Parallel()
	t.Run("", func(t *testing.T) {
		t.Parallel()
	})
}

func TestValidNoSubTest(t *testing.T) {
	t.Parallel()
}
