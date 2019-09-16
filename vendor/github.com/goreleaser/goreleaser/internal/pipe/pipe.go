// Package pipe provides generic erros for pipes to use.
package pipe

// ErrSnapshotEnabled happens when goreleaser is running in snapshot mode.
// It usually means that publishing and maybe some validations were skipped.
var ErrSnapshotEnabled = Skip("disabled during snapshot mode")

// ErrSkipPublishEnabled happens if --skip-publish is set.
// It means that the part of a Piper that publishes its artifacts was not run.
var ErrSkipPublishEnabled = Skip("publishing is disabled")

// ErrSkipSignEnabled happens if --skip-sign is set.
// It means that the part of a Piper that signs some things was not run.
var ErrSkipSignEnabled = Skip("artifact signing is disabled")

// ErrSkipValidateEnabled happens if --skip-validate is set.
// It means that the part of a Piper that validates some things was not run.
var ErrSkipValidateEnabled = Skip("validation is disabled")

// IsSkip returns true if the error is an ErrSkip
func IsSkip(err error) bool {
	_, ok := err.(ErrSkip)
	return ok
}

// ErrSkip occurs when a pipe is skipped for some reason
type ErrSkip struct {
	reason string
}

// Error implements the error interface. returns the reason the pipe was skipped
func (e ErrSkip) Error() string {
	return e.reason
}

// Skip skips this pipe with the given reason
func Skip(reason string) ErrSkip {
	return ErrSkip{reason: reason}
}
