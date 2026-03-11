package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG304 - potential file inclusion vulnerability
var SampleCodeG304 = []CodeSample{
	{[]string{`
package main

import (
"os"
"io/ioutil"
"log"
)

func main() {
	f := os.Getenv("tainted_file")
	body, err := ioutil.ReadFile(f)
	if err != nil {
	log.Printf("Error: %v\n", err)
	}
	log.Print(body)

}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
"os"
"log"
)

func main() {
	f := os.Getenv("tainted_file")
	body, err := os.ReadFile(f)
	if err != nil {
	log.Printf("Error: %v\n", err)
	}
	log.Print(body)

}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
  		title := r.URL.Query().Get("title")
		f, err := os.Open(title)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		body := make([]byte, 5)
		if _, err = f.Read(body); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Fprintf(w, "%s", body)
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
  		title := r.URL.Query().Get("title")
		f, err := os.OpenFile(title, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		body := make([]byte, 5)
		if _, err = f.Read(body); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Fprintf(w, "%s", body)
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"log"
	"os"
	"io/ioutil"
)

	func main() {
		f2 := os.Getenv("tainted_file2")
		body, err := ioutil.ReadFile("/tmp/" + f2)
		if err != nil {
		log.Printf("Error: %v\n", err)
	  }
		log.Print(body)
 }
 `}, 1, gosec.NewConfig()},
	{[]string{`
 package main

 import (
	 "bufio"
	 "fmt"
	 "os"
	 "path/filepath"
 )

func main() {
	reader := bufio.NewReader(os.Stdin)
  fmt.Print("Please enter file to read: ")
	file, _ := reader.ReadString('\n')
	file = file[:len(file)-1]
	f, err := os.Open(filepath.Join("/tmp/service/", file))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	contents := make([]byte, 15)
  if _, err = f.Read(contents); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
  fmt.Println(string(contents))
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
)

func main() {
	dir := os.Getenv("server_root")
	f3 := os.Getenv("tainted_file3")
	// edge case where both a binary expression and file Join are used.
	body, err := ioutil.ReadFile(filepath.Join("/var/"+dir, f3))
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Print(body)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
    "os"
    "path/filepath"
)

func main() {
    repoFile := "path_of_file"
    cleanRepoFile := filepath.Clean(repoFile)
    _, err := os.OpenFile(cleanRepoFile, os.O_RDONLY, 0600)
    if err != nil {
        panic(err)
    }
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
    "os"
    "path/filepath"
)

func openFile(filePath string) {
	_, err := os.OpenFile(filepath.Clean(filePath), os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
}

func main() {
    repoFile := "path_of_file"
	openFile(repoFile)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
    "os"
    "path/filepath"
)

func openFile(dir string, filePath string) {
	fp := filepath.Join(dir, filePath)
	fp = filepath.Clean(fp)
	_, err := os.OpenFile(fp, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
}

func main() {
    repoFile := "path_of_file"
	dir := "path_of_dir"
	openFile(dir, repoFile)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
    "os"
    "path/filepath"
)

func main() {
    repoFile := "path_of_file"
	relFile, err := filepath.Rel("./", repoFile)
	if err != nil {
		panic(err)
	}
    _, err = os.OpenFile(relFile, os.O_RDONLY, 0600)
    if err != nil {
        panic(err)
    }
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"io"
	"os"
)

func createFile(file string) *os.File {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	s, err := os.Open("src")
	if err != nil {
		panic(err)
	}
	defer s.Close()

	d := createFile("dst")
	defer d.Close()

	_, err = io.Copy(d, s)
	if  err != nil {
		panic(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"path/filepath"
)

type foo struct {
}

func (f *foo) doSomething(silly string) error {
	whoCares, err := filepath.Rel(THEWD, silly)
	if err != nil {
		return err
	}
	fmt.Printf("%s", whoCares)
	return nil
}

func main() {
	f := &foo{}

	if err := f.doSomething("irrelevant"); err != nil {
		panic(err)
	}
}
`, `
package main

var THEWD string
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"os"
	"path/filepath"
)

func open(fn string, perm os.FileMode) {
	fh, err := os.OpenFile(filepath.Clean(fn), os.O_RDONLY, perm)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
}

func main() {
	fn := "filename"
	open(fn, 0o600)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"os"
	"path/filepath"
)

func open(fn string, flag int) {
	fh, err := os.OpenFile(filepath.Clean(fn), flag, 0o600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
}

func main() {
	fn := "filename"
	open(fn, os.O_RDONLY)
}
`}, 0, gosec.NewConfig()},
}
