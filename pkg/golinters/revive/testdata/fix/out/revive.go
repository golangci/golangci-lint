//golangcitest:args -Erevive
//golangcitest:config_path testdata/revive-fix.yml
//golangcitest:expected_exitcode 0
package in

import (
	"errors"
	"fmt"
	"math"
)

func _() error {
	return fmt.Errorf("foo: %d", math.MaxInt)
}

func _() (int, error) {
	c := fmt.Errorf("bar: %d", math.MaxInt)
	return 1, c
}

func _() (int, error) {
	if c := fmt.Errorf("bar: %d", math.MaxInt); c != nil {
		return 0, c
	}
	return 1, nil
}
