package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG301 - mkdir permission check
var SampleCodeG301 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Mkdir("/tmp/mydir", 0777)
	if err != nil {
		fmt.Println("Error when creating a directory!")
		return
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.MkdirAll("/tmp/mydir", 0777)
	if err != nil {
		fmt.Println("Error when creating a directory!")
		return
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Mkdir("/tmp/mydir", 0600)
	if err != nil {
		fmt.Println("Error when creating a directory!")
		return
	}
}
`}, 0, gosec.NewConfig()},
}
