//args: -Erevive
//config_path: testdata/configs/revive.yml
package testdata

import "time"

func testRevive(t *time.Duration) error {
	if t == nil {
		return nil
	} else { // ERROR "if block ends with a return statement, so drop this else and outdent its block"
		return nil
	}
}
