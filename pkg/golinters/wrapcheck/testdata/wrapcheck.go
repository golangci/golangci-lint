//golangcitest:args -Ewrapcheck
package testdata

import (
	"encoding/json"
)

func do() error {
	_, err := json.Marshal(struct{}{})
	if err != nil {
		return err // want "error returned from external package is unwrapped"
	}

	return nil
}
