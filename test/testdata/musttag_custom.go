//golangcitest:args -Emusttag
//golangcitest:config_path testdata/configs/musttag.yml
package testdata

import (
	"encoding/asn1"
	"encoding/json"
)

// builtin functions:
func musttagJSONCustom() {
	var user struct { // want "`anonymous struct` should be annotated with the `json` tag as it is passed to `json.Marshal` at "
		Name  string
		Email string `json:"email"`
	}
	json.Marshal(user)
}

// custom functions from config:
func musttagASN1Custom() {
	var user struct { // want "`anonymous struct` should be annotated with the `asn1` tag as it is passed to `asn1.Marshal` at "
		Name  string
		Email string `asn1:"email"`
	}
	asn1.Marshal(user)
}
