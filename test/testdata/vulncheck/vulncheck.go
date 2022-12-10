package vulncheck

import "golang.org/x/text/language"

func testvuln() {
	_ = language.MustParseRegion("US")
}
