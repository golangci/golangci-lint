//golangcitest:args -Eforbidigo
//golangcitest:config_path testdata/configs/forbidigo_tests.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"time"
)

func Forbidigo() {
	fmt.Printf("too noisy!!!")
	time.Sleep(time.Nanosecond)
}
