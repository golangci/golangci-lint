package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG111 - potential directory traversal
var SampleCodeG111 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/bad/", http.StripPrefix("/bad/", http.FileServer(http.Dir("/"))))
	http.HandleFunc("/", HelloServer)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
`}, 1, gosec.NewConfig()},
}
