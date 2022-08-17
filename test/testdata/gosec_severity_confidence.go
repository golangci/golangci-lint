//golangcitest:args -Egosec
//golangcitest:config_path testdata/configs/gosec_severity_confidence.yml
package testdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var url string = "https://www.abcdefghijk.com"

func gosecVariableURL() {
	resp, err := http.Get(url) // ERROR "G107: Potential HTTP request made with variable url"
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", body)
}

func gosecHardcodedCredentials() {
	username := "admin"
	var password = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"

	fmt.Println("Doing something with: ", username, password)
}
