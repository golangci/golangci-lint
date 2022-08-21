//golangcitest:args -Einterfacer --internal-cmd-test
package testdata

import "io"

func InterfacerCheck(f io.ReadCloser) { // want "`f` can be `io.Closer`"
	f.Close()
}
