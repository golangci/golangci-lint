//golangcitest:args -Eerrgroupcheck
//golangcitest:config_path testdata/errgroupcheck_wait.yml
package testdata

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func ErrgroupWithWait() {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		return nil
	})

	_ = eg.Wait()
}

func ErrgroupMissingWait() {
	eg := errgroup.Group{} // want "errgroup 'eg' does not have Wait called"

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		return nil
	})
}

func ErrgroupContextWithWait() {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		return nil
	})

	_ = eg.Wait()
}

func ErrgroupContextMissingWait() {
	eg, _ := errgroup.WithContext(context.Background()) // want "errgroup 'eg' does not have Wait called"

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		return nil
	})
}

func ErrgroupMultipleScopesWithWait() {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		eg2 := errgroup.Group{}

		eg2.Go(func() error {
			return nil
		})

		return eg2.Wait()
	})

	_ = eg.Wait()
}

func ErrgroupMultipleScopesMissingWait() {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return nil
	})

	eg.Go(func() error {
		eg2 := errgroup.Group{} // want "errgroup 'eg2' does not have Wait called"

		eg2.Go(func() error {
			return nil
		})

		return nil
	})

	_ = eg.Wait()
}
