//args: -Egocritic
//config: linters-settings.gocritic.enabled-checks=ruleguard
//config: linters-settings.gocritic.settings.ruleguard.rules=ruleguard/rangeExprCopy.go
package p

func gocritic() {
	var xs [2048]byte

	// xs -> &xs
	for _, x := range xs {
		print(x)
	}
}
