//go:build debug

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

var (
	flagCPUProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	flagMemProfile = flag.String("memprofile", "", "write memory profile to file")
)

// Profiler manages CPU and memory profiling for debug builds.
// This encapsulation avoids package-level mutable state and enables proper testing.
type Profiler struct {
	cpuProfileFile *os.File
	logger         *log.Logger
	cleanupOnce    sync.Once
	cpuProfile     string
	memProfile     string
}

// NewProfiler creates a new profiler instance.
// If logger is nil, a no-op logger is used.
func NewProfiler(cpuProfile, memProfile string, logger *log.Logger) *Profiler {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &Profiler{
		cpuProfile: cpuProfile,
		memProfile: memProfile,
		logger:     logger,
	}
}

// Start begins CPU profiling if enabled.
// Returns an error if profiling cannot be started.
func (p *Profiler) Start() error {
	if p.cpuProfile == "" {
		return nil
	}

	f, err := os.Create(p.cpuProfile)
	if err != nil {
		return fmt.Errorf("could not create CPU profile: %w", err)
	}
	p.cpuProfileFile = f

	if err := pprof.StartCPUProfile(p.cpuProfileFile); err != nil {
		p.cpuProfileFile.Close()
		p.cpuProfileFile = nil
		return fmt.Errorf("could not start CPU profile: %w", err)
	}

	p.logger.Printf("CPU profiling enabled, writing to: %s", p.cpuProfile)
	return nil
}

// Stop writes memory profile and stops CPU profiling.
// Safe to call multiple times - only runs once.
// Logs errors but does not return them since this is cleanup code.
func (p *Profiler) Stop() {
	p.cleanupOnce.Do(func() {
		// Write memory profile
		if p.memProfile != "" {
			if err := p.writeMemoryProfile(); err != nil {
				p.logger.Printf("could not write memory profile: %v", err)
			} else {
				p.logger.Printf("memory profile written to: %s", p.memProfile)
			}
		}

		// Stop CPU profiling
		if p.cpuProfileFile != nil {
			pprof.StopCPUProfile()
			p.cpuProfileFile.Close()
			p.logger.Printf("CPU profile written to: %s", p.cpuProfile)
		}
	})
}

// writeMemoryProfile writes the memory profile to the configured file.
func (p *Profiler) writeMemoryProfile() error {
	f, err := os.Create(p.memProfile)
	if err != nil {
		return err
	}
	defer f.Close()

	runtime.GC() // get up-to-date statistics
	return pprof.WriteHeapProfile(f)
}

// initProfiling creates and starts profiling based on command-line flags.
// Returns the profiler instance and any error encountered during startup.
func initProfiling(logger *log.Logger) (*Profiler, error) {
	profiler := NewProfiler(*flagCPUProfile, *flagMemProfile, logger)
	if err := profiler.Start(); err != nil {
		return nil, err
	}
	return profiler, nil
}

// finishProfiling stops the profiler if it's not nil.
func finishProfiling(profiler *Profiler) {
	if profiler != nil {
		profiler.Stop()
	}
}
