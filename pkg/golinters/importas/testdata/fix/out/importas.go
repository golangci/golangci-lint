//golangcitest:args -Eimportas
//golangcitest:config_path testdata/importas.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	fff "fmt"
	"os"
	std_os "os"

	ananas "golang.org/x/tools/go/analysis"
)

func ImportAsWrongAlias() {
	fff.Println("foo")
	std_os.Stdout.WriteString("bar")
	os.Stdout.WriteString("test")
	_ = ananas.Analyzer{}
}
