//args: -Esamealias
package testdata

import (
	"fmt"
	alias1 "fmt"
	alias2 "fmt" // ERROR `package "fmt" have different alias, "alias2", "alias1"`
)

func SamealiasTest() {
	fmt.Println("test")
	alias1.Println("test")
	alias2.Println("test")
}
