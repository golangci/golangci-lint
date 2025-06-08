//golangcitest:args -Eptrequality
//golangcitest:config_path testdata/ptrequality.yml
package main

import (
	"errors"
	"log"
	"net/url"
)

func main() {
	_, err := url.Parse("://example.com")

	if errors.Is(err, &url.Error{}) { // want "is always false"
		log.Fatal("Cannot parse URL")
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		log.Fatalf("Cannot parse URL: %v", urlErr)
	}
}
