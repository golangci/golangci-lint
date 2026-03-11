package a

import (
	"crypto/md5"  // want `G501`
	"crypto/sha1" // want `G505`
	"os/exec"
)

func NosecVariants() {
	// Basic #nosec comment suppresses the diagnostic
	h1 := md5.New() // #nosec
	_ = h1

	// #nosec with rule ID
	h2 := md5.New() // #nosec G401
	_ = h2

	// #nosec with multiple rule IDs
	h3 := sha1.New() // #nosec G401 G505
	_ = h3

	// nosec without # should NOT suppress
	h4 := md5.New() // nosec // want `G401`
	_ = h4

	// Wrong rule ID should NOT suppress (G204 != G401)
	h5 := md5.New() // #nosec G204 -- wrong rule // want `G401`
	_ = h5

	// #nosec with explanation
	h6 := md5.New() // #nosec G401 -- used for non-cryptographic checksum
	_ = h6

	// Command injection with #nosec
	cmd := exec.Command("sh", "-c", getUserInput()) // #nosec G204
	_ = cmd
}
