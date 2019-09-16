package retry

// Error defines package errors.
type Error string

// Error implements the error interface.
func (err Error) Error() string {
	return string(err)
}

// Interrupted is the error returned by retry when the context is canceled.
const Interrupted Error = "operation interrupted"

// IsInterrupted checks that the error is related to the Breaker interruption.
// Deprecated: use err == retry.Interrupted instead.
// TODO:v5 will be removed
func IsInterrupted(err error) bool {
	return err == Interrupted
}

// IsRecovered checks that the error is related to unhandled Action's panic
// and returns an original cause of panic.
// Deprecated: retry won't handle unexpected panic and throw it.
// TODO:v5 will be removed
func IsRecovered(err error) (interface{}, bool) {
	if r, is := err.(result); is && r.recovered != nil {
		return r.recovered, true
	}
	return nil, false
}

type result struct {
	error
	recovered interface{}
}
