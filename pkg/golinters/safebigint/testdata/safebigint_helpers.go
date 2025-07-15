package safebigint

import (
	"math/big"
)

func unsupportedExpr() {
	// will be parsed but won't match selector expression unless it's in a call
	_ = big.NewInt(1).Add(big.NewInt(1), big.NewInt(2))
	_ = 123 // Still might not run unless passed to call
}

func literalArg(x *big.Int) {
	x.Add(big.NewInt(1), nil) // nil → getReferencedObject(nil) → default case
}

func noUses() {
	var x *big.Int
	x.BitLen()     // x is declared, but `TypesInfo.Uses[x]` may be nil
	_ = x.BitLen() // x used but not as selector.X — this may resolve to nil in TypesInfo.Uses
}

func callWithCompositeExpr() {
	// This expression is a *big.Int but not an Ident or SelectorExpr
	_ = []*big.Int{big.NewInt(1)}[0].Uint64()
}

type myInt int

func checkLocalNamed() {
	var x *myInt
	_ = x
}

func notASelector() {
	f := func() {}
	f() // call.Fun is an Ident, not a SelectorExpr
}
