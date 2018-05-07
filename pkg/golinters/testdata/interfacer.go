package testdata

import "io"

func InterfacerCheck(f io.ReadCloser) { // ERROR "XXX"
	f.Close()
}
