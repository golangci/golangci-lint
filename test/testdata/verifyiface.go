//args: -Everifyiface
package testdata

type Iface interface {
	Do() error
}

type NotChecked interface {
	X()
}

type Ok struct {
}

func (o Ok) Do() error {
	return nil
}

func (o Ok) X() {
}

type Fail struct { // ERROR `struct Fail doesn't verify interface compliance for testdata.Iface`
}

func (o *Fail) Do() error {
	return nil
}

func DoAny(x interface{}) {
	if iface, ok := x.(Iface); ok {
		iface.Do()
	}
}

var _ Iface = (*Ok)(nil)
