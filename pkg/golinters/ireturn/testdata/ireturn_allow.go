//golangcitest:args -Eireturn
//golangcitest:config_path testdata/ireturn_allow.yml
//golangcitest:expected_exitcode 0
package testdata

type (
	IreturnAllowDoer interface{ Do() }
	ireturnAllowDoer struct{}
)

func NewAllowDoer() IreturnAllowDoer { return new(ireturnAllowDoer) }
func (d *ireturnAllowDoer) Do()      { /*...*/ }

func NewerAllowDoer() *ireturnAllowDoer { return new(ireturnAllowDoer) }
