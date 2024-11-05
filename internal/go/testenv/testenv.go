// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testenv provides information about what functionality
// is available in different testing environments run by the Go team.
//
// It is an internal package because these details are specific
// to the Go team's test setup (on build.golang.org) and not
// fundamental to tests in general.
package testenv

// SyscallIsNotSupported reports whether err may indicate that a system call is
// not supported by the current platform or execution environment.
func SyscallIsNotSupported(err error) bool {
	return syscallIsNotSupported(err)
}
