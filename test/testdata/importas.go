//args: -Eimportas
//config_path: testdata/configs/importas.yml
package testdata

import (
	wrong_alias "fmt"      // ERROR `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	wrong_alias_again "os" // ERROR `import "os" imported as "wrong_alias_again" but must be "std_os" according to config`
)

func ImportAsWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
}
