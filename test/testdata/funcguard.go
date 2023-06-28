//golangcitest:args -Efuncguard
//golangcitest:config_path testdata/configs/funcguard.yml
package testdata

import "fmt"

func testFuncGuardDefault() {
	fmt.Println("hello") // want "do not use Println"
	return
}
