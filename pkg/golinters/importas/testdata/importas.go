//golangcitest:args -Eimportas
//golangcitest:config_path testdata/importas.yml
package testdata

import (
	wrong_alias "fmt" // want `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	"os"
	wrong_alias_again "os" // want `import "os" imported as "wrong_alias_again" but must be "std_os" according to config`

	wrong "golang.org/x/tools/go/analysis" // want `import "golang.org/x/tools/go/analysis" imported as "wrong" but must be "ananas" according to config`
)

func ImportAsWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = wrong.Analyzer{}
}
