//golangcitest:args -Exmlencoderclose
package testdata

import (
	"bytes"
	"encoding/xml"
)

func xmlEncoderClose() (string, error) {
	type document struct {
		A string `xml:"a"`
	}

	var buf bytes.Buffer
	err := xml.NewEncoder(&buf).Encode(document{ // want "Encoder.Close must be called"
		A: "abc123",
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
