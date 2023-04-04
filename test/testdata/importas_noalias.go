//golangcitest:args -Eimportas
//golangcitest:config_path testdata/configs/importas_noalias.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	wrong_alias "fmt"
	"os"
	wrong_alias_again "os"
)

func ImportAsNoAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
}
