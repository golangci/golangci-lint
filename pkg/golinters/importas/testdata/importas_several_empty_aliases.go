//golangcitest:args -Eimportas
//golangcitest:config_path testdata/importas_several_empty_aliases.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"math"
	"os"
)

func _() {
	fmt.Println("a")
	fmt.Fprint(os.Stderr, "b")
	println(math.MaxInt)
}
