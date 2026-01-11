//golangcitest:args -Etfproviderlint
//golangcitest:config_path tfproviderlint_no_extra.yml
package testdata

import (
	"fmt"
	"time"
)

func exampleResourceCreateNoExtra() {
	// With enable-extra: false, only standard passes run (not xpasses)
	// R018 is a standard pass, so it should still be detected
	time.Sleep(time.Second) // want "R018: prefer resource.Retry\\(\\) or \\(resource.StateChangeConf\\).WaitForState\\(\\) over time.Sleep\\(\\)"

	fmt.Println("resource created")
}
