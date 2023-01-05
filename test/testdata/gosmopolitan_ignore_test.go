//golangcitest:args -Egosmopolitan
//golangcitest:expected_exitcode 0
package testdata

import (
	"time"
)

func main() {
	_ = "默认不检查测试文件"
	_ = time.Local
}
