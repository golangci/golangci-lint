//golangcitest:args -Emnd
package testdata

import (
	"log"
	"net/http"
	"time"
)

func UseMagicNumber() {
	c := &http.Client{
		Timeout: 2 * time.Second, // want "Magic number: 2, in <assign> detected"
	}

	res, err := c.Get("http://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 { // want "Magic number: 200, in <condition> detected"
		log.Println("Something went wrong")
	}
}

func UseNoMagicNumber() {
	c := &http.Client{
		Timeout: time.Second,
	}

	res, err := c.Get("http://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Something went wrong")
	}
}
