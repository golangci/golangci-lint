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
