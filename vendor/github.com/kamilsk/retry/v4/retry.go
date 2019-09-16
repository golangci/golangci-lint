// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
package retry

import (
	"context"
	"sync/atomic"
)

// Retry takes action and performs it, repetitively, until successful.
// When it is done it releases resources associated with the Breaker.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(
	breaker BreakCloser,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	err := retry(breaker, action, strategies...)
	breaker.Close()
	return err
}

// Try takes action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Try(
	breaker Breaker,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	return retry(breaker, action, strategies...)
}

// TryContext takes action and performs it, repetitively, until successful.
// It uses the Context as a Breaker to prevent unnecessary action execution.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func TryContext(
	ctx context.Context,
	action func(ctx context.Context, attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	cascade, cancel := context.WithCancel(ctx)
	err := retry(ctx, currying(cascade, action), strategies...)
	cancel()
	return err
}

func currying(ctx context.Context, action func(context.Context, uint) error) func(uint) error {
	return func(attempt uint) error { return action(ctx, attempt) }
}

func retry(
	breaker Breaker,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	var interrupted uint32
	done := make(chan result, 1)

	go func(breaker *uint32) {
		var err error

		defer func() {
			done <- result{err, recover()}
			close(done)
		}()

		for attempt := uint(0); shouldAttempt(breaker, attempt, err, strategies...); attempt++ {
			err = action(attempt)
		}
	}(&interrupted)

	select {
	case <-breaker.Done():
		atomic.CompareAndSwapUint32(&interrupted, 0, 1)
		return Interrupted
	case err := <-done:
		if _, is := IsRecovered(err); is {
			return err
			// TODO:v5 throw origin
			// panic(origin)
		}
		return err.error
	}
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(breaker *uint32, attempt uint, err error, strategies ...func(uint, error) bool) bool {
	should := attempt == 0 || err != nil

	for i, repeat := 0, len(strategies); should && i < repeat; i++ {
		should = should && strategies[i](attempt, err)
	}

	return should && !atomic.CompareAndSwapUint32(breaker, 1, 0)
}
