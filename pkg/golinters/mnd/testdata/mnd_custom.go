//golangcitest:args -Emnd
//golangcitest:config_path testdata/mnd_custom.yml
package testdata

import (
	"log"
	"net/http"
	"time"
)

func Mnd() {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := c.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 { // want "Magic number: 200, in <condition> detected"
		log.Println("Something went wrong")
	}
}
