//golangcitest:args -Ereassign
package testdata

import "io"

func breakIO() {
	io.EOF = nil // ERROR `reassigning variable EOF in other package io`
}
