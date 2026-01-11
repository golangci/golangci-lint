//golangcitest:args -Etfproviderlint
package testdata

import (
	"fmt"
	"time"
)

func exampleResourceCreate() {
	// R018: time.Sleep should not be used in Terraform providers
	time.Sleep(time.Second) // want "R018: prefer resource.Retry\\(\\) or \\(resource.StateChangeConf\\).WaitForState\\(\\) over time.Sleep\\(\\)"

	fmt.Println("resource created")
}

func exampleResourceRead() {
	// Another time.Sleep violation
	time.Sleep(5 * time.Second) // want "R018: prefer resource.Retry\\(\\) or \\(resource.StateChangeConf\\).WaitForState\\(\\) over time.Sleep\\(\\)"
}
