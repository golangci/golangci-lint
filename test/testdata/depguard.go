package testdata

import (
	"log" // ERROR "`log` is in the blacklist"
)

func SpewDebugInfo() {
	log.Println("Debug info")
}
