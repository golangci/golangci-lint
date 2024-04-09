//golangcitest:args -Egosmopolitan
//golangcitest:config_path testdata/gosmopolitan_dont_ignore_tests.yml
package testdata

import (
	"time"
)

func main() {
	_ = "开启检查测试文件" // want `string literal contains rune in Han script`
	_ = time.Local // want `usage of time.Local`
}
