package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeG101 code snippets for hardcoded credentials
	SampleCodeG101 = []CodeSample{
		{[]string{`
package main

import "fmt"

func main() {
	username := "admin"
	password := "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
	fmt.Println("Doing something with: ", username, password)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
// Entropy check should not report this error by default
package main

import "fmt"

func main() {
	username := "admin"
	password := "secret"
	fmt.Println("Doing something with: ", username, password)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

var password = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"

func main() {
	username := "admin"
	fmt.Println("Doing something with: ", username, password)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const password = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"

func main() {
	username := "admin"
	fmt.Println("Doing something with: ", username, password)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	username = "user"
	password = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
)

func main() {
	fmt.Println("Doing something with: ", username, password)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

var password string

func init() {
	password = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

const (
	ATNStateSomethingElse = 1
	ATNStateTokenStart = 42
)

func main() {
	println(ATNStateTokenStart)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

const (
	ATNStateTokenStart = "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
)

func main() {
	println(ATNStateTokenStart)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var password string
	if password == "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" {
		fmt.Println("password equality")
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var password string
	if "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" == password {
		fmt.Println("password equality")
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var password string
	if password != "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" {
		fmt.Println("password equality")
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var password string
	if "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" != password {
		fmt.Println("password equality")
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var p string
	if p != "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" {
		fmt.Println("password equality")
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	var p string
	if "f62e5bcda4fae4f82370da0c6f20697b8f8447ef" != p {
		fmt.Println("password equality")
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	pw = "KjasdlkjapoIKLlka98098sdf012U/rL2sLdBqOHQUlt5Z6kCgKGDyCFA=="
)

func main() {
	fmt.Println(pw)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

var (
	pw string
)

func main() {
	pw = "KjasdlkjapoIKLlka98098sdf012U/rL2sLdBqOHQUlt5Z6kCgKGDyCFA=="
	fmt.Println(pw)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	cred = "KjasdlkjapoIKLlka98098sdf012U/rL2sLdBqOHQUlt5Z6kCgKGDyCFA=="
)

func main() {
	fmt.Println(cred)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

var (
	cred string
)

func main() {
	cred = "KjasdlkjapoIKLlka98098sdf012U/rL2sLdBqOHQUlt5Z6kCgKGDyCFA=="
	fmt.Println(cred)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	apiKey = "KjasdlkjapoIKLlka98098sdf012U"
)

func main() {
	fmt.Println(apiKey)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

var (
	apiKey string
)

func main() {
	apiKey = "KjasdlkjapoIKLlka98098sdf012U"
	fmt.Println(apiKey)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	bearer = "Bearer: 2lkjdfoiuwer092834kjdwf09"
)

func main() {
	fmt.Println(bearer)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

var (
	bearer string
)

func main() {
	bearer = "Bearer: 2lkjdfoiuwer092834kjdwf09"
	fmt.Println(bearer)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

// #nosec G101
const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config" // #nosec G101
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}

`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

// #nosec G101
const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config"
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}
	
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config" // #nosec G101
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}
	
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

//gosec:disable G101
const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config" //gosec:disable G101
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}

`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

//gosec:disable G101
const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config"
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}
	
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

const (
	ConfigLearnerTokenAuth string = "learner_auth_token_config" //gosec:disable G101
)

func main() {
	fmt.Printf("%s\n", ConfigLearnerTokenAuth)
}

`}, 0, gosec.NewConfig()},
		{[]string{`
package main

type DBConfig struct {
	Password string
}

func main() {
	_ = DBConfig{
		Password: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

type DBConfig struct {
	Password string
}

func main() {
	_ = &DBConfig{
		Password: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = struct{ Password string }{
		Password: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = map[string]string{
		"password": "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = map[string]string{
		"apiKey": "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

type Config struct {
	Username string
	Password string
}

func main() {
	_ = Config{
		Username: "admin",
		Password: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

type DBConfig struct {
	Password string
}

func main() {
	_ = DBConfig{ // #nosec G101
		Password: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef",
	}
}
`}, 0, gosec.NewConfig()},
		// Negatives
		{[]string{`
package main

func main() {
	_ = struct{ Password string }{
		Password: "secret", // low entropy
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = map[string]string{
		"password": "secret", // low entropy
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = struct{ Username string }{
		Username: "f62e5bcda4fae4f82370da0c6f20697b8f8447ef", // non-sensitive key
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = []string{"f62e5bcda4fae4f82370da0c6f20697b8f8447ef"} // unkeyed – no trigger
}
`}, 0, gosec.NewConfig()},
	}

	// SampleCodeG101Values code snippets for hardcoded credentials
	SampleCodeG101Values = []CodeSample{
		{[]string{`
package main

import "fmt"

func main() {
	customerNameEnvKey := "FOO_CUSTOMER_NAME"
	fmt.Println(customerNameEnvKey)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	txnID := "3637cfcc1eec55a50f78a7c435914583ccbc75a21dec9a0e94dfa077647146d7"
	fmt.Println(txnID)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	url := "https://username:abcdef0123456789abcdef0123456789abcdef01@contoso.com/"
	fmt.Println(url)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	githubToken := "ghp_iR54dhCYg9Tfmoywi9xLmmKZrrnAw438BYh3"
	fmt.Println(githubToken)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	awsAccessKeyID := "AKIAI44QH8DHBEXAMPLE"
	fmt.Println(awsAccessKeyID)
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	compareGoogleAPI := "test"
	if compareGoogleAPI == "AIzajtGS_aJGkoiAmSbXzu9I-1eytAi9Lrlh-vT" {
		fmt.Println(compareGoogleAPI)
	}
}	
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = struct{ SomeKey string }{
		SomeKey: "AKIAI44QH8DHBEXAMPLE",
	}
}
`}, 1, gosec.NewConfig()},
		{[]string{`
package main

func main() {
	_ = map[string]string{
		"github_token": "ghp_iR54dhCYg9Tfmoywi9xLmmKZrrnAw438BYh3",
	}
}
`}, 1, gosec.NewConfig()},
	}
)
