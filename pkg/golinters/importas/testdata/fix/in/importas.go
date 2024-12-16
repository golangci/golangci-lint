//golangcitest:args -Eimportas
//golangcitest:config_path testdata/importas.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	wrong_alias "fmt"
	"os"
	wrong_alias_again "os"

	wrong "golang.org/x/tools/go/analysis"
)

func ImportAsWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = wrong.Analyzer{}
}
