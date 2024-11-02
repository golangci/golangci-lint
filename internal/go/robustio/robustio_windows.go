// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package robustio

import (
	"errors"
	"syscall"
)

const errFileNotFound = syscall.ERROR_FILE_NOT_FOUND

// ERROR_SHARING_VIOLATION (ldez) extract from go1.19.1/src/internal/syscall/windows/syscall_windows.go.
// This is the only modification of this file.
const ERROR_SHARING_VIOLATION syscall.Errno = 32

// isEphemeralError returns true if err may be resolved by waiting.
func isEphemeralError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case syscall.ERROR_ACCESS_DENIED,
			syscall.ERROR_FILE_NOT_FOUND,
			ERROR_SHARING_VIOLATION:
			return true
		}
	}
	return false
}
