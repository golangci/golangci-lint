// (c) Copyright gosec's authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyzers

import (
	"testing"
)

// TestTaintAnalyzerConstructors tests that all taint analyzer constructors work.
func TestTaintAnalyzerConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor AnalyzerBuilder
		id          string
		description string
	}{
		{
			name:        "SQLInjection",
			constructor: newSQLInjectionAnalyzer,
			id:          "G701",
			description: "SQL injection via taint analysis",
		},
		{
			name:        "CommandInjection",
			constructor: newCommandInjectionAnalyzer,
			id:          "G702",
			description: "Command injection via taint analysis",
		},
		{
			name:        "PathTraversal",
			constructor: newPathTraversalAnalyzer,
			id:          "G703",
			description: "Path traversal via taint analysis",
		},
		{
			name:        "SSRF",
			constructor: newSSRFAnalyzer,
			id:          "G704",
			description: "SSRF via taint analysis",
		},
		{
			name:        "XSS",
			constructor: newXSSAnalyzer,
			id:          "G705",
			description: "XSS via taint analysis",
		},
		{
			name:        "LogInjection",
			constructor: newLogInjectionAnalyzer,
			id:          "G706",
			description: "Log injection via taint analysis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := tt.constructor(tt.id, tt.description)

			if analyzer == nil {
				t.Fatal("constructor returned nil")
			}

			if analyzer.Name != tt.id {
				t.Errorf("analyzer Name = %s, want %s", analyzer.Name, tt.id)
			}

			if analyzer.Run == nil {
				t.Error("analyzer Run function is nil")
			}

			if len(analyzer.Requires) == 0 {
				t.Error("analyzer has no requirements")
			}
		})
	}
}

// TestDefaultAnalyzersIncludeTaint tests that default analyzers include taint rules.
func TestDefaultAnalyzersIncludeTaint(t *testing.T) {
	expectedTaintIDs := []string{"G701", "G702", "G703", "G704", "G705", "G706"}

	found := make(map[string]bool)
	for _, def := range defaultAnalyzers {
		found[def.ID] = true
	}

	for _, id := range expectedTaintIDs {
		if !found[id] {
			t.Errorf("default analyzers missing taint rule: %s", id)
		}
	}
}

// TestGenerateIncludesTaintAnalyzers tests that Generate includes taint analyzers.
func TestGenerateIncludesTaintAnalyzers(t *testing.T) {
	analyzerList := Generate(false)

	expectedTaintIDs := []string{"G701", "G702", "G703", "G704", "G705", "G706"}

	for _, id := range expectedTaintIDs {
		if _, ok := analyzerList.Analyzers[id]; !ok {
			t.Errorf("generated analyzer list missing taint rule: %s", id)
		}
	}
}

// TestGenerateExcludeTaintAnalyzers tests that taint analyzers can be excluded.
func TestGenerateExcludeTaintAnalyzers(t *testing.T) {
	filter := NewAnalyzerFilter(true, "G701", "G702")
	analyzerList := Generate(false, filter)

	if _, ok := analyzerList.Analyzers["G701"]; ok {
		t.Error("G701 should be excluded but was found")
	}

	if _, ok := analyzerList.Analyzers["G702"]; ok {
		t.Error("G702 should be excluded but was found")
	}

	// Other taint analyzers should still be present
	if _, ok := analyzerList.Analyzers["G703"]; !ok {
		t.Error("G703 should be present but was not found")
	}
}
