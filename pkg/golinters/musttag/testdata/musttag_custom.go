//golangcitest:args -Emusttag
//golangcitest:config_path testdata/musttag.yml
package testdata

import (
	"encoding/asn1"
	"encoding/json"
)

// builtin functions:
func musttagJSONCustom() {
	var user struct {
		Name  string
		Email string `json:"email"`
	}
	json.Marshal(user) // want "the given struct should be annotated with the `json` tag"
}

// custom functions from config:
func musttagASN1Custom() {
	var user struct {
		Name  string
		Email string `asn1:"email"`
	}
	asn1.Marshal(user) // want "the given struct should be annotated with the `asn1` tag"
}
