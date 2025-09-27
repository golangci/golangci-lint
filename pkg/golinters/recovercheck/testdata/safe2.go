package testdata

import aliaspkg "recovercheck/pkg"

// SafeGoroutineWithAliasImport uses a recovery function from another package with import alias
func SafeGoroutineWithAliasImport() {
	go func() {
		defer aliaspkg.PanicRecover()
		panic("oh no")
	}()
}
