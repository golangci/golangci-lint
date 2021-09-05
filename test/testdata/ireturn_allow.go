// args: -Eireturn
// config: linters-settings.ireturn.allow=["IreturnAllowDoer"]
package testdata

type (
	IreturnAllowDoer interface{ Do() }
	ireturnAllowDoer struct{}
)

func NewAllowDoer() IreturnAllowDoer { return new(ireturnAllowDoer) }
func (d *ireturnAllowDoer) Do()      { /*...*/ }

func NewerAllowDoer() *ireturnAllowDoer { return new(ireturnAllowDoer) }
