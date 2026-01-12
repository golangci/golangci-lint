//golangcitest:config_path testdata/tfproviderlint_custom.yml
//golangcitest:args -Etfproviderlint
package tfproviderlint

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func _() {
	var d schema.ResourceData

	d.HasChanges("attr1", "attr2", "attr3", "attr4", "attr5", "attr6") // want `R019: d\.HasChanges\(\) has many arguments, consider d\.HasChangesExcept\(\)`
	d.HasChanges("attr1", "attr2", "attr3")                            // want `R019: d\.HasChanges\(\) has many arguments, consider d\.HasChangesExcept\(\)`
	d.HasChanges("attr1", "attr2")
}

func _() {
	time.Sleep(time.Second) // want `R018: prefer resource.Retry\(\) or \(resource.StateChangeConf\).WaitForState\(\) over time.Sleep\(\)`

	fmt.Println("resource created")
}

func _() {
	// Another time.Sleep violation
	time.Sleep(5 * time.Second) // want `R018: prefer resource.Retry\(\) or \(resource.StateChangeConf\).WaitForState\(\) over time.Sleep\(\)`
}

func _() {
	var d schema.ResourceData

	d.GetOkExists("test")
}
