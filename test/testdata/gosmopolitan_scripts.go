//golangcitest:args -Egosmopolitan
//golangcitest:config_path testdata/configs/gosmopolitan_scripts.yml
package testdata

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")                 // want `string literal contains rune in Latin script`
	fmt.Println("should not report this line") //nolint:gosmopolitan
	fmt.Println("你好，世界")
	fmt.Println("こんにちは、セカイ") // want `string literal contains rune in Hiragana script`
}
