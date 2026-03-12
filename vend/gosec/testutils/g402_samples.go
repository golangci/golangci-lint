package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG402 - TLS settings
var SampleCodeG402 = []CodeSample{
	{[]string{`
// InsecureSkipVerify
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// InsecureSkipVerify from variable
package main

import (
	"crypto/tls"
)

func main() {
	var conf tls.Config
	conf.InsecureSkipVerify = true
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{MinVersion: 0},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main

import (
	"crypto/tls"
	"fmt"
)

func CaseNotError() *tls.Config {
	var v uint16 = tls.VersionTLS13

	return &tls.Config{
		MinVersion: v,
	}
}

func main() {
    a := CaseNotError()
	fmt.Printf("Debug: %v\n", a.MinVersion)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main

import (
	"crypto/tls"
	"fmt"
)

func CaseNotError() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
}

func main() {
    a := CaseNotError()
	fmt.Printf("Debug: %v\n", a.MinVersion)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main
import (
	"crypto/tls"
	"fmt"
)

func CaseError() *tls.Config {
	var v = &tls.Config{
		MinVersion: 0,
	}
	return v
}

func main() {
    a := CaseError()
	fmt.Printf("Debug: %v\n", a.MinVersion)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main

import (
	"crypto/tls"
	"fmt"
)

func CaseError() *tls.Config {
	var v = &tls.Config{
		MinVersion: getVersion(),
	}
	return v
}

func getVersion() uint16 {
	return tls.VersionTLS12
}

func main() {
    a := CaseError()
	fmt.Printf("Debug: %v\n", a.MinVersion)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Insecure minimum version
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

var theValue uint16 = 0x0304

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{MinVersion: theValue},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Insecure max version
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{MaxVersion: 0},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// Insecure ciphersuite selection
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// secure max version when min version is specified
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: 0,
			MinVersion: tls.VersionTLS13,
		},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package p0

import "crypto/tls"

func TlsConfig0() *tls.Config {
	var v uint16 = 0
	return &tls.Config{MinVersion: v}
}
`, `
package p0

import "crypto/tls"

func TlsConfig1() *tls.Config {
   return &tls.Config{MinVersion: 0x0304}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	cfg := tls.Config{
		MinVersion: MinVer,
	}
	fmt.Println("tls min version", cfg.MinVersion)
}
`, `
package main

import "crypto/tls"

const MinVer = tls.VersionTLS13
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"crypto/tls"
	cryptotls "crypto/tls"
)

func main() {
	_ = tls.Config{MinVersion: tls.VersionTLS12}
	_ = cryptotls.Config{MinVersion: cryptotls.VersionTLS12}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// InsecureSkipVerify with unary NOT (direct !false → true, high confidence)
package main

import "crypto/tls"

func main() {
	_ = &tls.Config{InsecureSkipVerify: !false}
}
`}, 1, gosec.NewConfig()},

	{[]string{`
// InsecureSkipVerify with unary NOT (direct !true → false, no issue)
package main

import "crypto/tls"

func main() {
	_ = &tls.Config{InsecureSkipVerify: !true}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// InsecureSkipVerify via const with NOT (resolves to true, high confidence)
package main

import "crypto/tls"

const skipVerify = !false

func main() {
	_ = &tls.Config{InsecureSkipVerify: skipVerify}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// PreferServerCipherSuites false (direct, medium severity)
package main

import "crypto/tls"

func main() {
	_ = &tls.Config{PreferServerCipherSuites: false}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// PreferServerCipherSuites with !true (resolves to false)
package main

import "crypto/tls"

func main() {
	_ = &tls.Config{PreferServerCipherSuites: !true}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// PreferServerCipherSuites true (no issue)
package main

import "crypto/tls"

func main() {
	_ = &tls.Config{PreferServerCipherSuites: true}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// MaxVersion explicitly low via variable
package main

import "crypto/tls"

func main() {
	var lowMax uint16 = tls.VersionTLS10
	_ = &tls.Config{MaxVersion: lowMax}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// PreferServerCipherSuites unknown → low-confidence
package main

import "crypto/tls"

var prefer bool // unresolved

func main() {
	_ = &tls.Config{PreferServerCipherSuites: prefer}
}
`}, 1, gosec.NewConfig()},
}
