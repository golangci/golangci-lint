//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet_fix.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"log"
	"os"
)

type Foo struct {
	C string
	A []string
	E int32
	B bool
	D int8
}

func nonConstantFormat(s string) {
	fmt.Printf("%s", s)
	fmt.Printf(s, "arg")
	fmt.Fprintf(os.Stderr, "%s", s)
	log.Printf("%s", s)
}
