//args: -Emultiimport
package testdata

import (
	"fmt"      // ERROR "import appears multiple times under different aliases"
	blah "fmt" // ERROR "import appears multiple times under different aliases"
)

func MultiImportIssue() {
	fmt.Println(`a`)
	blah.Println(`a`)
}
