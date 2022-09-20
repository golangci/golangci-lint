//golangcitest:args -Etodolint
//golangcitest:config_path testdata/configs/todolint.yml
package testdata

import "fmt"

func todoLintWithOptionsExample() {
	// TODO: This is ignored due to keywords config

	fmt.Println("你好，世界") // XXX: more languages // want `TODO comment should be in the form XXX\(author\)`
}
