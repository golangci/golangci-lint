//golangcitest:args -Emusttag
package testdata

import (
	"encoding/asn1"
	"encoding/json"
)

// builtin functions:
func musttagJSON() {
	var user struct { // want "`anonymous struct` should be annotated with the `json` tag as it is passed to `json.Marshal` at "
		Name  string
		Email string `json:"email"`
	}
	json.Marshal(user)
}

// custom functions from config:
func musttagASN1() {
	var user struct {
		Name  string
		Email string `asn1:"email"`
	}
	asn1.Marshal(user)
}
