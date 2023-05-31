//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/configs/forbidigo_struct.yml
package testdata

import (
	fmt2 "fmt"
	"time"
)

func Forbidigo() {
	fmt2.Printf("too noisy!!!") // want "use of `fmt2\\.Printf` forbidden by pattern `fmt\\\\.Print\\.\\*`"
	time.Sleep(time.Nanosecond) // want "no sleeping!"
}
