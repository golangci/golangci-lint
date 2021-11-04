//args: -Eexhaustive
//config_path: testdata/configs/exhaustive_checking_strategy_name.yml
package testdata

type AccessControl string

const (
	AccessPublic  AccessControl = "public"
	AccessPrivate AccessControl = "private"
	AccessDefault AccessControl = AccessPublic
)

func example(v AccessControl) {
	switch v { // ERROR "missing cases in switch of type AccessControl: AccessDefault"
	case AccessPublic:
	case AccessPrivate:
	}
}
