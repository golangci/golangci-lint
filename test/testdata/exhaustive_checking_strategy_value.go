//args: -Eexhaustive
//config_path: testdata/configs/exhaustive_checking_strategy_value.yml
package testdata

type AccessControl string

const (
	AccessPublic  AccessControl = "public"
	AccessPrivate AccessControl = "private"
	AccessDefault AccessControl = AccessPublic
)

// Expect no diagnostics for this switch statement, even though AccessDefault is
// not listed, because AccessPublic (which is listed) has the same value as
// AccessDefault.

func example(v AccessControl) {
	switch v {
	case AccessPublic:
	case AccessPrivate:
	}
}
