package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG107 - SSRF via http requests with variable url
var SampleCodeG107 = []CodeSample{
	{[]string{`
// Input from the std in is considered insecure
package main
import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"bufio"
)
func main() {
	in := bufio.NewReader(os.Stdin)
	url, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", body)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Variable defined a package level can be changed at any time
// regardless of the initial value
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var url string = "https://www.google.com"

func main() {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", body)
}`}, 1, gosec.NewConfig()},
	{[]string{`
// Environmental variables are not considered as secure source
package main
import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
)
func main() {
	url := os.Getenv("tainted_url")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
			panic(err)
	}
	fmt.Printf("%s", body)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Constant variables or hard-coded strings are secure
package main

import (
	"fmt"
	"net/http"
)
const url = "http://127.0.0.1"
func main() {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// A variable at function scope which is initialized to
// a constant string is secure (e.g. cannot be changed concurrently)
package main

import (
	"fmt"
	"net/http"
)
func main() {
	var url string = "http://127.0.0.1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// A variable at function scope which is initialized to
// a constant string is secure (e.g. cannot be changed concurrently)
package main

import (
	"fmt"
	"net/http"
)
func main() {
	url := "http://127.0.0.1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// A variable at function scope which is initialized to
// a constant string is secure (e.g. cannot be changed concurrently)
package main

import (
	"fmt"
	"net/http"
)
func main() {
	url1 := "test"
	var url2 string = "http://127.0.0.1"
	url2 = url1
	resp, err := http.Get(url2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// An exported variable declared a packaged scope is not secure
// because it can changed at any time
package main

import (
	"fmt"
	"net/http"
)

var Url string

func main() {
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// An url provided as a function argument is not secure
package main

import (
	"fmt"
	"net/http"
)
func get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}
func main() {
	url := "http://127.0.0.1"
	get(url)
}
`}, 1, gosec.NewConfig()},
}
