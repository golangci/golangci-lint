package a

import (
	"crypto/md5" // want `G501: \[CWE-327\] Blocklisted import crypto/md5: weak cryptographic primitive`
	"fmt"
	"os/exec"
)

func VulnerableFunction() {
	// Test SQL injection - gosec doesn't catch simple string concatenation without database/sql
	query := "SELECT * FROM users WHERE name = '" + getUserInput() + "'"
	_ = query

	// G204: Command injection (AST-based rule)
	cmd := exec.Command("sh", "-c", getUserInput()) // want `G204: \[CWE-78\] Subprocess launched with a potential tainted input or cmd arguments`
	_ = cmd

	// G401: Weak crypto (AST-based rule)
	h := md5.New() // want `G401: \[CWE-328\] Use of weak cryptographic primitive`
	_ = h
}

func getUserInput() string {
	return "test"
}

func SecureFunction() {
	fmt.Println("This is secure")
}

func IntegerOverflow() {
	// G115: Integer overflow in type conversion (SSA-based analyzer)
	var a uint32 = 0xFFFFFFFF
	b := int32(a) // want `G115`
	fmt.Println(b)
}

func SliceBounds() {
	// G602: Slice bounds check (SSA-based analyzer)
	s := []int{1, 2, 3}
	idx := 10
	_ = s[:idx] // want `G602`
}
