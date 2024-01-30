//go:build ruleguard

package ruleguard

import "github.com/quasilyte/go-ruleguard/dsl"

func preferWriteString(m dsl.Matcher) {
	m.Match(`$w.Write([]byte($s))`).
		Where(m["w"].Type.Implements("io.StringWriter")).
		Suggest("$w.WriteString($s)").
		Report(`$w.WriteString($s) should be preferred to the $$`)
}
