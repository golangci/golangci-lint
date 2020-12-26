// +build !go1.14

package test

import "testing"

const tCleanupExists = false

// do nothing for go.1.13
func registerCleanup(t *testing.T, f func()) {}
