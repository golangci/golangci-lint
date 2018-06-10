// args: -Edepguard --depguard.include-go-root --depguard.packages='log'
package testdata

import (
	"log" // ERROR "`log` is in the blacklist"
)

func SpewDebugInfo() {
	log.Println("Debug info")
}
