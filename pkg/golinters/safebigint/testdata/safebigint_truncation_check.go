package safebigint

import (
	"math/big"
)

type myBigInt struct{}

func (m *myBigInt) Uint64() uint64 { return 42 } // Should NOT trigger

func safeConversion(b *big.Int) uint64 {
	if b.Cmp(big.NewInt(0)) < 0 {
		return 0
	}
	return b.Uint64() // want "calling Uint64 on \\*big.Int may silently truncate or overflow"
}

func testMixed() {
	x := big.NewInt(123)
	_ = x.Uint64() // want "calling Uint64 on \\*big.Int may silently truncate or overflow"

	_ = x.BitLen() // not a truncating method

	y := new(myBigInt)
	_ = y.Uint64() // OK: user-defined type
}

func testIgnore() {
	var i int
	_ = uint64(i) // OK: not big.Int
}

func otherTruncationExamples() {
	b := big.NewInt(123456789)

	_ = b.Uint64() // want "calling Uint64 on \\*big.Int may silently truncate or overflow"
	_ = b.Int64()  // want "calling Int64 on \\*big.Int may silently truncate or overflow"
}
