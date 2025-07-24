//golangcitest:args -Etagliatelle
//golangcitest:config_path testdata/tagliatelle_initialism_overrides.yml
package testdata

type Foo struct {
	UserAMQP   string `json:"useAmqp"` // want `json\(camel\): got 'useAmqp' want 'useAMQP'`
	FooLHS     string `json:"fooLhs"`
	FooRHS     string `json:"fooRhs"`     // want `json\(camel\): got 'fooRhs' want 'fooRHS'`
	SomeAWSKey string `json:"someAwsKey"` // want `json\(camel\): got 'someAwsKey' want 'someAWSKey'`
}
