//args: -Elooppointer
package testdata

func lintLoopPointer() {
	type domain struct {
		property *int
	}

	type dto struct {
		prop int
		cond bool
	}

	result := make([]domain, 0, 3)
	for _, d := range []dto{
		{prop: 10},
		{prop: 11},
		{prop: 12},
	} {
		p := domain{}
		if !d.cond {
			p.property = &d.prop // ERROR "taking a pointer for the loop variable d"
		}
		result = append(result, p)
	}

	for _, r := range result {
		if r.property != nil {
			// expected 10, 11, 12
			// obtained 12, 12, 12
			println(*r.property)
		}
	}
}
