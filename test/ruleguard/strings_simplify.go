//go:build ruleguard

package ruleguard

import "github.com/quasilyte/go-ruleguard/dsl"

func StringsSimplify(m dsl.Matcher) {
	// Some issues have simple fixes that can be expressed as
	// a replacement pattern. Rules can use Suggest() function
	// to add a quickfix action for such issues.
	m.Match(`strings.Replace($s, $old, $new, -1)`).
		Report(`this Replace call can be simplified`).
		Suggest(`strings.ReplaceAll($s, $old, $new)`)

	// Suggest() can be used without Report().
	// It'll print the suggested template to the user.
	m.Match(`strings.Count($s1, $s2) == 0`).
		Suggest(`!strings.Contains($s1, $s2)`)
}
