//golangcitest:args -Eimportas
//golangcitest:config_path testdata/configs/importas_strict.yml
package testdata

import (
	wrong_alias "fmt"      // want `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	"os"                   // want `import "os" imported without alias but must be with alias "std_os" according to config`
	wrong_alias_again "os" // want `import "os" imported as "wrong_alias_again" but must be "std_os" according to config`

	wrong "github.com/pkg/errors" // want `import "github.com/pkg/errors" imported as "wrong" but must be "pkgerr" according to config`
)

func ImportAsStrictWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = wrong.New("baz")
}
