package testdata

import (
	aliaspkg "recovercheck/pkg"

	"golang.org/x/sync/errgroup"
)

func SafeErrgroupWithRecoverFromAliasImport() {
	var g errgroup.Group

	// This should NOT be flagged - uses recovery from another package with import alias
	g.Go(func() error {
		defer aliaspkg.PanicRecover()
		panic("This panic is recovered by another package with import alias")
		return nil
	})

	g.Wait()
}
