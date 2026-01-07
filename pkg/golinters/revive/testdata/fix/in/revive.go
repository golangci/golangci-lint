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
	return errors.New(fmt.Sprintf("foo: %d", math.MaxInt))
}

func _() (int, error) {
	c := errors.New(fmt.Sprintf("bar: %d", math.MaxInt))
	return 1, c
}

func _() (int, error) {
	if c := errors.New(fmt.Sprintf("bar: %d", math.MaxInt)); c != nil {
		return 0, c
	}
	return 1, nil
}
