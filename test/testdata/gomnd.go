//args: -Egomnd
package testdata

import (
	"log"
	"net/http"
	"time"
)

func UseMagicNumber() {
	c := &http.Client{
		Timeout: 1 * time.Second, // ERROR : "Magic number: 1, in <assign> detected"
	}

	res, err := c.Get("http://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 { // ERROR : "Magic number: 200, in <condition> detected"
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
