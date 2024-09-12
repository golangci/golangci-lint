//golangcitest:args -Ereceivernamelen
package testdata

type Example1 struct{}

func (e *Example1) F() {}

type Example2 struct{}

func (ex *Example2) F() {}

type Example3 struct{}

func (ex3 *Example3) F() {} // want `receiver variable names must be one or two letters in length`

type Example4 struct{}

func (example *Example4) F() {} // want `receiver variable names must be one or two letters in length`

func F() {}
