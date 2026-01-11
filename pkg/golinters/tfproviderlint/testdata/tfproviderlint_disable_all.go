//golangcitest:args -Etfproviderlint
//golangcitest:config_path tfproviderlint_disable_all.yml
package testdata

import (
	"fmt"
	"time"
)

func exampleResourceCreateDisableAll() {
	// With disable-all: true and enable: [R018], only R018 should run
	time.Sleep(time.Second) // want "R018: prefer resource.Retry\\(\\) or \\(resource.StateChangeConf\\).WaitForState\\(\\) over time.Sleep\\(\\)"

	fmt.Println("resource created")
}
