//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/forbidigo.yml
package testdata

import (
	"fmt"
	fmt2 "fmt"
	"time"
)

func Forbidigo() {
	fmt.Printf("too noisy!!!")  // want "use of `fmt\\.Printf` forbidden by pattern `fmt\\\\.Print\\.\\*`"
	fmt2.Printf("too noisy!!!") // Not detected because analyze-types is false by default for backward compatibility.
	time.Sleep(time.Nanosecond) // want "no sleeping!"
}
