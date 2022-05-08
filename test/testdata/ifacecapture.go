// args: -Eifacecapture
package testdata

// Example program that erroneously captures the outer variable when it likely
// intends to use the parameter interface.

type MyInterface interface {
	Do()
}

type MyImpl struct{}

var _ MyInterface = (*MyImpl)(nil)

func (m *MyImpl) Do() {}

func doThing(callback func(tx MyInterface)) {
	myImpl := MyImpl{}
	callback(&myImpl)
}

type HasMyImpl struct {
	A MyImpl
}

func (h HasMyImpl) GetMyImpl() *MyImpl {
	return &h.A
}

func main() {
	outer := MyImpl{}
	outer2 := HasMyImpl{A: MyImpl{}}
	outer3 := struct{ B HasMyImpl }{B: HasMyImpl{A: MyImpl{}}}
	outerArr := [2]MyImpl{{}, {}}
	doThing(func(inner MyInterface) {
		outer.Do()              // want "captured variable outer implements interface MyInterface"
		outer2.A.Do()           // want "captured variable outer2.A implements interface MyInterface"
		outer3.B.A.Do()         // want "captured variable outer3.B.A implements interface MyInterface"
		outerArr[0].Do()        // We don't flag this yet because it is a lot of extra work
		outer2.GetMyImpl().Do() // We don't flag this yet because it becomes much harder to analyze where the receiver is coming from
		inner.Do()
	})
}
