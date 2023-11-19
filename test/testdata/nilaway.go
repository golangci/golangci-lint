package testdata

type P struct {
	f int
}

// nilAwayExample1 tests NilAway's detection of uninitialized variable access.
func nilAwayExample1(someCondition bool) {
	var p *P
	if someCondition {
		p = &P{}
	}
	// NilAway should report a potential nil pointer dereference here if someCondition is false.
	if p != nil {
		print(p.f)
	}
}

// nilAwayExample2 tests NilAway's detection of nil returns across function boundaries.
func nilAwayExample2() {
	print(*nilAwayFoo()) // NilAway should report an error here.
}

func nilAwayFoo() *int {
	return nil
}
