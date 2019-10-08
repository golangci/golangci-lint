//args: -Egosec
package testdata

import (
	"crypto/md5" // ERROR "G501: Blacklisted import `crypto/md5`: weak cryptographic primitive"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Gosec() {
	h := md5.New() // ERROR "G401: Use of weak cryptographic primitive"
	log.Print(h)
}

func GosecNolintGas() {
	h := md5.New() //nolint:gas
	log.Print(h)
}

func GosecNolintGosec() {
	h := md5.New() //nolint:gosec
	log.Print(h)
}

func GosecNoErrorCheckingByDefault() {
	f, _ := os.Create("foo")
	fmt.Println(f)
}

func GosecG204SubprocWithFunc() {
	arg := func() string {
		return "/tmp/dummy"
	}

	exec.Command("ls", arg()).Run() // ERROR "G204: Subprocess launched with function call as argument or cmd arguments"
}
