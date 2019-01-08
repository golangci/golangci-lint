//args: -Egosimple
package testdata

import "strings"

func Gosimple(id1, s1 string) string {
	if strings.HasPrefix(id1, s1) { // ERROR "should replace.*with.*strings.TrimPrefix"
		id1 = strings.TrimPrefix(id1, s1)
	}
	return id1
}

func GosimpleNolintGosimple(id1, s1 string) string {
	if strings.HasPrefix(id1, s1) { //nolint:gosimple
		id1 = strings.TrimPrefix(id1, s1)
	}
	return id1
}

func GosimpleNolintMegacheck(id1, s1 string) string {
	if strings.HasPrefix(id1, s1) { //nolint:megacheck
		id1 = strings.TrimPrefix(id1, s1)
	}
	return id1
}
