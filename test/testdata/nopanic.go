//args: -Enopanic --internal-cmd-test
package testdata

import (
	"log"
	"os"
)

func NoPanic() {
	panic("message")     // ERROR `calling panic is not allowed`
	log.Panic("message") // ERROR `calling panic is not allowed`
}

func NoExit() {
	os.Exit(1)  // ERROR `program exit is not allowed`
	log.Fatal() // ERROR `program exit is not allowed`
}
