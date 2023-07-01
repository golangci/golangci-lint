//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/configs/path-except.yml
package testdata

import (
	"fmt"
	"testing"
	"time"
)

func TestForbidigo(t *testing.T) {
	fmt.Printf("too noisy!!!")  // want "use of `fmt\\.Printf` forbidden by pattern `fmt\\\\.Print\\.\\*`"
	time.Sleep(time.Nanosecond) // want "no sleeping!"
}
