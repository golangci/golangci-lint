//args: -Estructfield
//config: linters-settings.structfield.limit=0
package testdata

import (
	"fmt"
)

type Account struct {
	Name        string
	Email       string
	Permissions []Permission
	Verified    bool
	Deactivated bool
}

type Permission struct {
	Domain string
	Access string
}

func useStructLiteral() {
	acc := Account{ // ERROR "found 5 non-labeled fields on struct literal \\(> 0\\)"
		"John Smith",
		"john.smith@example.com",
		[]Permission{
			Permission{"account", "read"},  // ERROR "found 2 non-labeled fields on struct literal \\(> 0\\)"
			Permission{"account", "write"}, // ERROR "found 2 non-labeled fields on struct literal \\(> 0\\)"
		},
		true,
		false,
	}
	fmt.Printf("%+v", acc)
}
