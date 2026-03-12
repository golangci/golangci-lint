//go:build !debug

package main

import "log"

// Profiler is a stub type for release builds.
type Profiler struct{}

// initProfiling is a no-op in release builds.
// Profiling is only available when building with -tags debug.
func initProfiling(_ *log.Logger) (*Profiler, error) {
	return nil, nil
}

// finishProfiling is a no-op in release builds.
// Profiling is only available when building with -tags debug.
func finishProfiling(_ *Profiler) {}
