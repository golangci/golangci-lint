package testdata

func nilErr1() error {
	var p *P
	if someCondition {
		p = &P{}
	}
	print(p.f) // nilness reports NO error here, but NilAway does.
}

func nilErr2() error {
	func nilErr3() *int {
		return nil
	}

	print(*foo()) // nilness reports NO error here, but NilAway does.
}

func nilErr3() *int {
    return nil
}
