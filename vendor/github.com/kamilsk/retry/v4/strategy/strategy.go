// Package strategy provides a way to define how retry is performed.
package strategy

import "time"

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning `true`
// allows for the next attempt to be made. Returning `false` halts the retrying
// process and returns the last error returned by the called Action.
//
// The strategy will be passed an "attempt" number on each successive retry
// iteration, starting with a `0` value before the first attempt is actually
// made. This allows for a pre-action delay, etc.
type Strategy func(attempt uint, err error) bool

// Limit creates a Strategy that limits the number of attempts that Retry will
// make.
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint, _ error) bool {
		return attempt < attemptLimit
	}
}

// Delay creates a Strategy that waits the given duration before the first
// attempt is made.
func Delay(duration time.Duration) Strategy {
	return func(attempt uint, _ error) bool {
		if 0 == attempt {
			time.Sleep(duration)
		}

		return true
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(attempt uint, _ error) bool {
		if 0 < attempt && 0 < len(durations) {
			durationIndex := int(attempt - 1)

			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}

			time.Sleep(durations[durationIndex])
		}

		return true
	}
}

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given backoff.Algorithm.
func Backoff(algorithm func(attempt uint) time.Duration) Strategy {
	return BackoffWithJitter(algorithm, noJitter())
}

// BackoffWithJitter creates a Strategy that waits before each attempt, with a
// duration as defined by the given backoff.Algorithm and jitter.Transformation.
func BackoffWithJitter(
	algorithm func(attempt uint) time.Duration,
	transformation func(duration time.Duration) time.Duration,
) Strategy {
	return func(attempt uint, _ error) bool {
		if 0 < attempt {
			time.Sleep(transformation(algorithm(attempt)))
		}

		return true
	}
}

// noJitter creates a jitter.Transformation that simply returns the input.
func noJitter() func(time.Duration) time.Duration {
	return func(duration time.Duration) time.Duration {
		return duration
	}
}
