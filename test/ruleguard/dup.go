//go:build ruleguard

package ruleguard

import "github.com/quasilyte/go-ruleguard/dsl"

// Suppose that we want to report the duplicated left and right operands of binary operations.
//
// But if the operand has some side effects, this rule can cause false positives:
// `f() && f()` can make sense (although it's not the best piece of code).
//
// This is where *filters* come to the rescue.
func DupSubExpr(m dsl.Matcher) {
	// All filters are written as a Where() argument.
	// In our case, we need to assert that $x is "pure".
	// It can be achieved by checking the m["x"] member Pure field.
	m.Match(`$x || $x`,
		`$x && $x`,
		`$x | $x`,
		`$x & $x`).
		Where(m["x"].Pure).
		Report(`suspicious identical LHS and RHS`)
}
