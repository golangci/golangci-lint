package retry

// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of context.Context and github.com/kamilsk/breaker.Breaker.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
}

// A BreakCloser carries a cancellation signal to break an action execution
// and can release resources associated with it.
//
// It is a subset of github.com/kamilsk/breaker.Breaker.
type BreakCloser interface {
	Breaker
	// Close closes the Done channel and releases resources associated with it.
	Close()
}

// Action defines a callable function that package retry can handle.
type Action func(attempt uint) error

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How []func(attempt uint, err error) bool

// Interface defines a behavior of stateful executor of Actions in parallel.
// TODO:v5 complete the draft
type Interface interface {
	Try(Breaker, Action, ...How) Interface
}
