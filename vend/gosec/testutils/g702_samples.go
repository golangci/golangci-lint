package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG702 - Command injection via taint analysis
var SampleCodeG702 = []CodeSample{
	{[]string{`
package main

import (
	"net/http"
	"os/exec"
)

func handler(r *http.Request) {
	filename := r.URL.Query().Get("file")
	cmd := exec.Command("cat", filename)
	cmd.Run()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"os"
	"os/exec"
)

func dynamicCommand() {
	userInput := os.Args[1]
	exec.Command("sh", "-c", userInput).Run()
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"os/exec"
)

func safeCommand() {
	// Safe - no user input
	exec.Command("ls", "-la").Run()
}
`}, 0, gosec.NewConfig()},
}
