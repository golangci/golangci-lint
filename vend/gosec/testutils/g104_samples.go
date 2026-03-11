package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeG104 finds errors that aren't being handled
	SampleCodeG104 = []CodeSample{
		{[]string{`
package main

import "fmt"

func test() (int,error) {
	return 0, nil
}

func main() {
	v, _ := test()
	fmt.Println(v)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import (
	"io/ioutil"
	"os"
	"fmt"
)

func a() error {
	return fmt.Errorf("This is an error")
}

func b() {
	fmt.Println("b")
	ioutil.WriteFile("foo.txt", []byte("bar"), os.ModeExclusive)
}

func c() string {
	return fmt.Sprintf("This isn't anything")
}

func main() {
	_ = a()
	a()
	b()
	c()
}
`}, 2, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func test() error {
	return nil
}

func main() {
	e := test()
	fmt.Println(e)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
// +build go1.10

package main

import "strings"

func main() {
	var buf strings.Builder
	_, err := buf.WriteString("test string")
	if err != nil {
		panic(err)
	}
}`, `
package main

func dummy(){}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import (
	"bytes"
)

type a struct {
	buf *bytes.Buffer
}

func main() {
	a := &a{
		buf: new(bytes.Buffer),
	}
	a.buf.Write([]byte{0})
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import (
	"io/ioutil"
	"os"
	"fmt"
)

func a() {
	fmt.Println("a")
	ioutil.WriteFile("foo.txt", []byte("bar"), os.ModeExclusive)
}

func main() {
	a()
}
`}, 0, gosec.Config{"G104": map[string]interface{}{"ioutil": []interface{}{"WriteFile"}}}},
		{[]string{`
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func createBuffer() *bytes.Buffer {
	return new(bytes.Buffer)
}

func main() {
	new(bytes.Buffer).WriteString("*bytes.Buffer")
	fmt.Fprintln(os.Stderr, "fmt")
	new(strings.Builder).WriteString("*strings.Builder")
	_, pw := io.Pipe()
	pw.CloseWithError(io.EOF)

	createBuffer().WriteString("*bytes.Buffer")
	b := createBuffer()
	b.WriteString("*bytes.Buffer")
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "crypto/rand"

func main() {
	b := make([]byte, 8)
	rand.Read(b)
	_ = b
}
`}, 0, gosec.NewConfig()},
	} // it shouldn't return any errors because all method calls are whitelisted by default

	// SampleCodeG104Audit finds errors that aren't being handled in audit mode
	SampleCodeG104Audit = []CodeSample{
		{[]string{`
package main

import "fmt"

func test() (int,error) {
	return 0, nil
}

func main() {
	v, _ := test()
	fmt.Println(v)
}
`}, 1, gosec.Config{gosec.Globals: map[gosec.GlobalOption]string{gosec.Audit: "enabled"}}},
		{[]string{`
package main

import (
	"io/ioutil"
	"os"
	"fmt"
)

func a() error {
	return fmt.Errorf("This is an error")
}

func b() {
	fmt.Println("b")
	ioutil.WriteFile("foo.txt", []byte("bar"), os.ModeExclusive)
}

func c() string {
	return fmt.Sprintf("This isn't anything")
}

func main() {
	_ = a()
	a()
	b()
	c()
}
`}, 3, gosec.Config{gosec.Globals: map[gosec.GlobalOption]string{gosec.Audit: "enabled"}}},
		{[]string{`
package main

import "fmt"

func test() error {
	return nil
}

func main() {
	e := test()
	fmt.Println(e)
}
`}, 0, gosec.Config{gosec.Globals: map[gosec.GlobalOption]string{gosec.Audit: "enabled"}}},
		{[]string{`
// +build go1.10

package main

import "strings"

func main() {
	var buf strings.Builder
	_, err := buf.WriteString("test string")
	if err != nil {
		panic(err)
	}
}
`, `
package main

func dummy(){}
`}, 0, gosec.Config{gosec.Globals: map[gosec.GlobalOption]string{gosec.Audit: "enabled"}}},
	}
)
