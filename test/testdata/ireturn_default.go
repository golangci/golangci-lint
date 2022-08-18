//golangcitest:args -Eireturn
package testdata

type (
	IreturnDoer interface{ Do() }
	ireturnDoer struct{}
)

func New() IreturnDoer     { return new(ireturnDoer) } // ERROR `New returns interface \(command-line-arguments.IreturnDoer\)`
func (d *ireturnDoer) Do() { /*...*/ }

func Newer() *ireturnDoer { return new(ireturnDoer) }
