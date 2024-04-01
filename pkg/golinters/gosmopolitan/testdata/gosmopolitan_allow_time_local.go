//golangcitest:args -Egosmopolitan
//golangcitest:config_path testdata/gosmopolitan_allow_time_local.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"time"
)

func main() {
	_ = time.Local
}
