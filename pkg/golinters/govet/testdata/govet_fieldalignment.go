//go:build !(386 || arm || mips || mipsle)

//golangcitest:args -Egovet
//golangcitest:config_path testdata/govet_fieldalignment.yml
package testdata

type gvfaGood struct {
	y int32
	x byte
	z byte
}

type gvfaBad struct { // want `gvfaBad has size 12 \(allocator size class 16\) but the optimal size is 8 leading to a waste of 8 bytes`
	x byte
	y int32
	z byte
}

type gvfaPointerGood struct {
	P   *int
	buf [1000]uintptr
}

type gvfaPointerBad struct { // want "gvfaPointerBad has 8008 leading bytes of pointer data but optimal value is 8"
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

type gvfaPointerSortaBad struct { // want "gvfaPointerSortaBad has 32 leading bytes of pointer data but optimal value is 24"
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

type gvfaZeroBad struct { // want "gvfaZeroBad has size 8 but the optimal size is 4"
	a uint32
	b [0]byte
}
