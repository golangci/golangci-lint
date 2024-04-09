//golangcitest:args -Emusttag
package testdata

import (
	"encoding/asn1"
	"encoding/json"
)

// builtin functions:
func musttagJSON() {
	var user struct {
		Name  string
		Email string `json:"email"`
	}
	json.Marshal(user) // want "the given struct should be annotated with the `json` tag"
}

// custom functions from config:
func musttagASN1() {
	var user struct {
		Name  string
		Email string `asn1:"email"`
	}
	asn1.Marshal(user)
}
