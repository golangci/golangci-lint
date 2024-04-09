//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet_fieldalignment.yml
package testdata

type gvfaGood struct {
	y int32
	x byte
	z byte
}

type gvfaBad struct { // want "struct of size 12 could be 8"
	x byte
	y int32
	z byte
}

type gvfaPointerGood struct {
	P   *int
	buf [1000]uintptr
}

type gvfaPointerBad struct { // want "struct with 8008 pointer bytes could be 8"
	buf [1000]uintptr
	P   *int
}

type gvfaPointerSorta struct {
	a struct {
		p *int
		q uintptr
	}
	b struct {
		p *int
		q [2]uintptr
	}
}

type gvfaPointerSortaBad struct { // want "struct with 32 pointer bytes could be 24"
	a struct {
		p *int
		q [2]uintptr
	}
	b struct {
		p *int
		q uintptr
	}
}

type gvfaZeroGood struct {
	a [0]byte
	b uint32
}

type gvfaZeroBad struct { // want "struct of size 8 could be 4"
	a uint32
	b [0]byte
}
