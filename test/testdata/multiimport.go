//args: -Emultiimport
package testdata

import (
	"fmt"      // ERROR "import appears multiple times under different aliases"
	blah "fmt" // ERROR "import appears multiple times under different aliases"
	blah2 "log"
)

func MultiImportIssue() {
	fmt.Println(`a`)
	blah.Println(`a`)
	blah2.Print(`hi`)
}
