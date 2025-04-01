//golangcitest:args -Egocritic
//golangcitest:config_path testdata/gocritic_inportShadow.yml
package testdata

import (
	"fmt"
	"path/filepath"
)

func Bar() {
	filepath.Join("a", "b")
}

func foo() {
	filepath := "foo.txt" // want "importShadow: shadow of imported package 'filepath'"
	fmt.Printf("File path: %s\n", filepath)
}
