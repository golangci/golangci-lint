// Package ctrlc provides an easy way of having a task that is
// context-aware and that also deals with interrupt and term signals.
package ctrlc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Task is function that can be executed by a ctrlc
type Task func() error

// Ctrlc is the task ctrlc
type Ctrlc struct {
	signals chan os.Signal
	errs    chan error
}

// New returns a new ctrlc with its internals setup.
func New() *Ctrlc {
	return &Ctrlc{
		signals: make(chan os.Signal, 1),
		errs:    make(chan error, 1),
	}
}

// Default ctrlc instance
var Default = New()

// Run executes a given task with a given context, dealing with its timeouts,
// cancels and SIGTERM and SIGINT signals.
// It will return an error if the context is canceled, if deadline exceeds,
// if a SIGTERM or SIGINT is received and of course if the task itself fails.
func (c *Ctrlc) Run(ctx context.Context, task Task) error {
	go func() {
		c.errs <- task()
	}()
	signal.Notify(c.signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-c.errs:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-c.signals:
		return fmt.Errorf("received: %s", sig)
	}
}
