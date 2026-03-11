package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG112 - potential slowloris attack
var SampleCodeG112 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	err := (&http.Server{
		Addr: ":1234",
	}).ListenAndServe()
	if err != nil {
		panic(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"time"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	server := &http.Server{
		Addr:              ":1234",
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"time"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	server := &http.Server{
		Addr:              ":1234",
		ReadTimeout:  	   1 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	hs  *http.Server
	mux *http.ServeMux
	mu  sync.Mutex
}

func New(listenAddr string) *Server {
	mux := http.NewServeMux()

	return &Server{
	hs: &http.Server{ // #nosec G112 - Not publicly exposed
		Addr:    listenAddr,
		Handler: mux,
	},
	mux: mux,
	mu:  sync.Mutex{},
	}
}

func main() {
	fmt.Print("test")
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	hs  *http.Server
	mux *http.ServeMux
	mu  sync.Mutex
}

func New(listenAddr string) *Server {
	mux := http.NewServeMux()

	return &Server{
	hs: &http.Server{ //gosec:disable G112 - Not publicly exposed
		Addr:    listenAddr,
		Handler: mux,
	},
	mux: mux,
	mu:  sync.Mutex{},
	}
}

func main() {
	fmt.Print("test")
}
`}, 0, gosec.NewConfig()},
}
