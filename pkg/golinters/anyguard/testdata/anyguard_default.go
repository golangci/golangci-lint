//golangcitest:args -Eanyguard
//golangcitest:expected_exitcode 0
package testdata

import "fmt"

var _ = fmt.Sprintf

type DefaultAllowedPayload map[string]any
