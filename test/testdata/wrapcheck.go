//args: -Ewrapcheck
package testdata

import (
	"encoding/json"
)

func do() error {
	_, err := json.Marshal(struct{}{})
	if err != nil {
		return err // ERROR "error returned from external package is unwrapped"
	}

	return nil
}
