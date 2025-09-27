// Package errgroup provides a mock implementation for testing purposes
package errgroup

// Group is a mock errgroup.Group for testing
type Group struct{}

// Go runs the given function in a new goroutine
func (g *Group) Go(f func() error) {
	go f()
}

// Wait waits for all goroutines to complete
func (g *Group) Wait() error {
	return nil
}
