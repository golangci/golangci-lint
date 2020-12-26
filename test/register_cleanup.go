// +build go1.14

package test

import "testing"

const tCleanupExists = true

// registerCleanup exists because t.Cleanup doesn't exist prior to go1.14
func registerCleanup(t *testing.T, f func()) {
	t.Cleanup(f)
}
