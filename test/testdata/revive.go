//args: -Erevive
//config_path: testdata/configs/revive.yml
package testdata

import "time"

func testRevive(t *time.Duration) error {
	if t == nil {
		return nil
	} else { // ERROR "indent-error-flow: if block ends with a return statement, .*"
		return nil
	}
}
