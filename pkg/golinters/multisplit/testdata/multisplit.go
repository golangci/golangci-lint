//golangcitest:args -Emultisplit
package ungroup

func value() int {
	return 1
}

func value2() (int, int) {
	return 1, 2
}

// var declaration without initializer

var vdBad1, vdBad2 int // want `variable declaration with multiple identifiers \(vdBad1, vdBad2\) should be split into individual declarations`

var (
	vdBad3, vdBad4 int // want `variable declaration with multiple identifiers \(vdBad3, vdBad4\) should be split into individual declarations`
	vdOther        string
)

func vd() {
	var vdBad5, vdBad6 int

	_ = vdBad5
	_ = vdBad6
}

// =>
var (
	vdGood1 int
	vdGood2 int
)

// var declaration initializer

var vdiBad1, vdiBad2 = 1, 2 // want `variable declaration with multiple identifiers and initializers \(vdiBad1, vdiBad2\) should be split into individual declarations`

func vdi() {
	// no short declaration on purpose
	var vdiBad7, vdiBad8 = 1, 2

	_ = vdiBad7
	_ = vdiBad8
}

// =>
var (
	vdiGood1 = 1
	vdiGood2 = 2
)

// typed

var vditBad1, vditBad2 int = 1, 2 // want `variable declaration with multiple identifiers and initializers \(vditBad1, vditBad2\) should be split into individual declarations`

var (
	vditBad3, vditBad4 int    = 1, 2 // want `variable declaration with multiple identifiers and initializers \(vditBad3, vditBad4\) should be split into individual declarations`
	vditOther          string = "three"
)

var vdiBad9, vdiBad10 = 1, value() // want `variable declaration with multiple identifiers and initializers \(vdiBad9, vdiBad10\) should be split into individual declarations`

func vdiOtherFn() {
	// no short declaration on purpose
	var vdiBad11, vdiBad12 = 1, value()

	_ = vdiBad11
	_ = vdiBad12
}

func vdit() {
	var vditBad5, vditBad6 int = 1, 2

	_ = vditBad5
	_ = vditBad6
}

// =>
var (
	vditGood1 int = 1
	vditGood2 int = 2
)

var (
	vditGood5 int = 1
	vditGood6 int = value()
)

// const declarations

const cdBad1, cdBad2 = 1, 2 // want `const declaration with multiple identifiers \(cdBad1, cdBad2\) should be split into individual declarations`

const (
	cdBad3, cdBad4 = 1, 2 // want `const declaration with multiple identifiers \(cdBad3, cdBad4\) should be split into individual declarations`
)

func cd() {
	const cdBad5, cdBad6 = 1, 2
}

// =>
const (
	cdGood1 = 1
	cdGood2 = 2
)

// typed

const cdtBad1, cdtBad2 int = 1, 2 // want `const declaration with multiple identifiers \(cdtBad1, cdtBad2\) should be split into individual declarations`

const (
	cdtBad3, cdtBad4 int = 1, 2 // want `const declaration with multiple identifiers \(cdtBad3, cdtBad4\) should be split into individual declarations`
)

func cdt() {
	const cdtBad5, cdtBad6 int = 1, 2
}

// =>
const (
	cdtGood1 int = 1
	cdtGood2 int = 2
)

// function parameters

func fBad1(p1, p2 int)            {} // want `function parameters with multiple identifiers \(p1, p2\) should be split into individual parameters`
func fBad2(p3, p4 int, p5 string) {} // want `function parameters with multiple identifiers \(p3, p4\) should be split into individual parameters`

type (
	fBadT1 func(p6, p7 int) string // want `function parameters with multiple identifiers \(p6, p7\) should be split into individual parameters`
	fBadT2 = func(p8, p9 int) string // want `function parameters with multiple identifiers \(p8, p9\) should be split into individual parameters`
)

func fBadM(fn func(p10, p11 int) string) {} // want `function parameters with multiple identifiers \(p10, p11\) should be split into individual parameters`

type iBad interface {
	m(p12, p13 int) string // want `function parameters with multiple identifiers \(p12, p13\) should be split into individual parameters`
}

// =>
func fGood(p14 int, p15 int) {}

// function return values

func frBad() (r1, r2 int) { // want `function return values with multiple identifiers \(r1, r2\) should be split into individual return values`
	return 0, 0
}

func frBad2() (r3, r4 int, r5 string) { // want `function return values with multiple identifiers \(r3, r4\) should be split into individual return values`
	return 0, 0, "three"
}

type frBadT1 func() (r6, r7 int, r8 string) // want `function return values with multiple identifiers \(r6, r7\) should be split into individual return values`

type frBadT2 = func() (r9, r10 int, r11 string) // want `function return values with multiple identifiers \(r9, r10\) should be split into individual return values`

func frBadM() func() (r10 string, r12, r13 int, r14 string) { // want `function return values with multiple identifiers \(r12, r13\) should be split into individual return values`
	return func() (r15 string, r16, r17 int, r18 string) { // want `function return values with multiple identifiers \(r16, r17\) should be split into individual return values`
		return "zero", 0, 0, "three"
	}
}

type iFrBad interface {
	m() (r19, r20 int, r21 string) // want `function return values with multiple identifiers \(r19, r20\) should be split into individual return values`
}

// =>
func frGood() (r22 int, r23 int) {
	return 0, 0
}

// assignments

var (
	aBad1  int
	aBad2  int
	aBad3  string
	aGood1 int
	aGood2 int
	aGood3 string
)

func aBad() {
	aBad1, aBad2, aBad3 = 1, 2, "three"
}

// =>
func aGood() {
	aGood1 = 1
	aGood2 = 2
	aGood3 = "three"
}

// not split able

func aUnsplitable() {
	for i, j := 0, 10; i < 10; i, j = i+1, j-1 {
	}

	aGood1, aGood2 = value2()
}

// short variable declarations

func svd1() {
	svdBad1, svdBad2 := 1, 2
	_, _ = svdBad1, svdBad2
}

// =>
func svd2() {
	svdGood1 := 1
	svdGood2 := 2
	_, _ = svdGood1, svdGood2
}

// struct fields

type SBad struct {
	val1, val2 int // want `struct field declaration with multiple identifiers \(val1, val2\) should be split into individual fields`
}

func sBad() {
	type tmp struct {
		val3, val4 int // want `struct field declaration with multiple identifiers \(val3, val4\) should be split into individual fields`
	}
}

// =>
type SGood struct {
	val1 int
	val2 int
}
