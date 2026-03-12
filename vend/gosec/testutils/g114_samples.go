package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG114 - Use of net/http serve functions that have no support for setting timeouts
var SampleCodeG114 = []CodeSample{
	{[]string{`
package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = http.Serve(l, nil)
	log.Fatal(err)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	l, err := net.Listen("tcp", ":8443")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = http.ServeTLS(l, nil, "cert.pem", "key.pem")
	log.Fatal(err)
}
`}, 1, gosec.NewConfig()},
}
