package iface

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
)

func TestAnalyzersFromSettings(t *testing.T) {
	testCases := map[string]struct {
		enable          []string
		expectedEnabled []string
	}{
		"nil analyzers": {
			enable:          nil,
			expectedEnabled: []string{"unused", "empty", "duplicate", "opaque"},
		},
		"empty analyzers": {
			enable:          []string{},
			expectedEnabled: []string{"unused", "empty", "duplicate", "opaque"},
		},
		"unused only": {
			enable:          []string{"unused"},
			expectedEnabled: []string{"unused"},
		},
		"some analyzers": {
			enable:          []string{"unused", "opaque"},
			expectedEnabled: []string{"unused", "opaque"},
		},
		"duplicate analyzers": {
			enable:          []string{"unused", "opaque", "unused"},
			expectedEnabled: []string{"unused", "opaque"},
		},
		"all analyzers": {
			enable:          []string{"unused", "opaque", "empty", "duplicate"},
			expectedEnabled: []string{"unused", "empty", "duplicate", "opaque"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			settings := &config.IfaceSettings{Enable: tc.enable}
			analyzers := analyzersFromSettings(settings)

			if len(analyzers) != len(tc.expectedEnabled) {
				t.Errorf("expected %d analyzers, got %d", len(tc.enable), len(analyzers))
			}

		LoopSettings:
			for _, a := range analyzers {
				for _, name := range tc.expectedEnabled {
					if a.Name == name {
						continue LoopSettings
					}
				}

				t.Errorf("unexpected analyzer %q", a.Name)
			}
		})
	}
}
