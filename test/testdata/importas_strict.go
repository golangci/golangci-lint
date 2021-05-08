//args: -Eimportas
//config_path: testdata/configs/importas_strict.yml
package testdata

import (
	wrong_alias "fmt"      // ERROR `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	"os"                   // ERROR `import "os" imported without alias but must be with alias "std_os" according to config`
	wrong_alias_again "os" // ERROR `import "os" imported as "wrong_alias_again" but must be "std_os" according to config`

	wrong "github.com/pkg/errors" // ERROR `import "github.com/pkg/errors" imported as "wrong" but must be "pkgerr" according to config`
)

func ImportAsStrictWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = wrong.New("baz")
}
