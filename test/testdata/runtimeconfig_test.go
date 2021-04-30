// args: -Eunused
package testdata

// TESTS WORK LOCALLY BUT NOT ON GITHUB DUE TO GRAPHENE BEING IN A PRIVATE REPO AND NOT ACCESSIBLE
// REENABLE ONCE THE LINTER GAINS ACCESS, DEFINE RUNTIMECONFIG LINTER VIA `// args: -Eruntimeconfig`
// CHANGE ALL __ERROR__ TO ERROR

/*
import (
	"ghe.anduril.dev/anduril/graphene-go/pkg/graphene"
)

// valid json tags -----------------------------------
type T1 struct {
	A string `json:"a"`
}
type T2 struct {
	A string `json:"b"`
}
// valid json tags -----------------------------------

// invalid json tags -----------------------------------
type T3 struct {
	A string `json:"a_b"`
}
type T4 struct {
	A string `json:"a-b"`
}
type T5 struct {
	A string
}
// invalid json tags -----------------------------------

func GetTest() {
	rtc := setup()

	// nil or primitive values are not allowed
	rtc.Get("", nil, nil)    // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
	rtc.Get("", T1{}, nil)   // __ERROR__ "expected configuration object as 3rd arg but got nil"
	rtc.Get("", nil, &T1{})  // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
	rtc.Get("", true, false) // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
	rtc.Get("", "", "")      // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
	rtc.Get("", 123, 456)    // __ERROR__ "expected 2nd arg to be non nil but got <nil>"

	// 3rd arg needs to be a ref (update in-place)
	t1 := T1{}
	rtc.Get("", T1{}, t1) // __ERROR__ "expected ref as 3rd arg to Get"

	// type mismatch
	t2 := &T2{} // __ERROR__ "the configuration object \\(arg t2\\) has to match the type of the default argument, expected type: command-line-arguments.T1 but found T2"
	rtc.Get("", T1{}, t2)

	// misconfigured tags
	t3 := &T3{}
	rtc.Get("", T3{}, t3) // __ERROR__ `runtimeConfigurations must specify a json field name in camelCase format \(e.g. json:'fieldName'\) for exported fields, found field 'A' of 'command-line-arguments.T3' using 'json:"a_b"'`
	t4 := &T4{}
	rtc.Get("", T4{}, t4) // __ERROR__ `runtimeConfigurations must specify a json field name in camelCase format \(e.g. json:'fieldName'\) for exported fields, found field 'A' of 'command-line-arguments.T4' using 'json:"a-b"'`
	t5 := &T5{}
	rtc.Get("", T5{}, t5) // __ERROR__ "runtimeConfigurations must specify a json field name in camelCase format \\(e.g. json:'fieldName'\\) for exported fields, found field 'A' of 'command-line-arguments.T5' using ''"
	// misconfigured tags

	// valid usage
	rtc.Get("", T1{}, &t1)
}

func SubscribeTest() {
	rtc := setup()

	// nil args
	rtc.Subscribe("", nil, nil)                                       // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
	rtc.Subscribe("", nil, func(c interface{}) *graphene.ApplyError { // __ERROR__ "expected 2nd arg to be non nil but got <nil>"
		return nil
	})
	t1 := T1{}
	rtc.Subscribe("", t1, nil) // __ERROR__ "expected function as 3rd arg but got nil"
	// nil args

	// type conversion requires ref
	rtc.Subscribe("", t1, func(c interface{}) *graphene.ApplyError {
		dummy(c.(T1)) // __ERROR__ "the configuration object \\(arg c\\) is a reference, add '\\*' to fix type conversion"
		return nil
	})

	// mismatched types and no ref
	rtc.Subscribe("", t1, func(c interface{}) *graphene.ApplyError {
		dummy(c.(T2)) // __ERROR__ "the configuration object \\(arg c\\) is a reference, add '\\*' to fix type conversion"
		return nil
	})

	// mismatched types
	rtc.Subscribe("", t1, func(c interface{}) *graphene.ApplyError {
		dummy(c.(*T2)) // __ERROR__ "the configuration object \\(arg c\\) has to match the type of the default argument, expected type: command-line-arguments.T1 but found T2"
		return nil
	})

	// misconfigured tags
	t3 := T3{}
	rtc.Subscribe("", t3, nil) // __ERROR__ `runtimeConfigurations must specify a json field name in camelCase format \(e.g. json:'fieldName'\) for exported fields, found field 'A' of 'command-line-arguments.T3' using 'json:"a_b"'`
	t4 := T4{}
	rtc.Subscribe("", t4, nil) // __ERROR__ `runtimeConfigurations must specify a json field name in camelCase format \(e.g. json:'fieldName'\) for exported fields, found field 'A' of 'command-line-arguments.T4' using 'json:"a-b"'`
	t5 := T5{}
	rtc.Subscribe("", t5, nil) // __ERROR__ "runtimeConfigurations must specify a json field name in camelCase format \\(e.g. json:'fieldName'\\) for exported fields, found field 'A' of 'command-line-arguments.T5' using ''"
	// misconfigured tags

	// valid usage
	rtc.Subscribe("", t1, func(c interface{}) *graphene.ApplyError {
		dummy(c.(*T1))
		return nil
	})
}

func dummy(_ interface{}) {}

func setup() graphene.RuntimeConfigRegistry {
	var rtc graphene.RuntimeConfigRegistry
	graphene.Create("asdf", 123, 123, nil, func(g *graphene.Graphene) error {
		rtc = g.RuntimeConfig()
		return nil
	})
	return rtc
}

*/
