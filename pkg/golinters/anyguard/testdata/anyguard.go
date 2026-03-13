//golangcitest:args -Eanyguard
//golangcitest:config_path testdata/anyguard.yml
package testdata

import "fmt"

var _ = fmt.Sprintf

type AllowedPayload map[string]any

type ViolatingPayload map[string]any // want "disallowed any usage"
