//golangcitest:args -Etfproviderlint
//golangcitest:config_path tfproviderlint_disable.yml
//golangcitest:expected_exitcode 0
package testdata

import (
	"fmt"
	"time"
)

func exampleResourceCreateDisabled() {
	// R018 is disabled, so this should NOT be flagged
	time.Sleep(time.Second)

	fmt.Println("resource created")
}

func exampleResourceReadDisabled() {
	// R018 is disabled, so this should NOT be flagged
	time.Sleep(5 * time.Second)
}
