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
	"go/token"
	"testing"

	"golang.org/x/tools/go/ssa"
)

// TestExtractBinOpBound_NilGuards tests nil safety in extractBinOpBound
func TestExtractBinOpBound_NilGuards(t *testing.T) {
	// Test nil binop
	bound, value, err := extractBinOpBound(nil)
	if err == nil {
		t.Error("expected error for nil binop")
	}
	if bound != lowerUnbounded {
		t.Errorf("expected lowerUnbounded, got %v", bound)
	}
	if value != 0 {
		t.Errorf("expected value 0, got %d", value)
	}
}

// TestExtractLenBound_NilGuards tests nil safety in extractLenBound
func TestExtractLenBound_NilGuards(t *testing.T) {
	// Test nil binop
	val, offset, ok := extractLenBound(nil)
	if ok {
		t.Error("expected false for nil binop")
	}
	if val != nil {
		t.Errorf("expected nil value, got %v", val)
	}
	if offset != 0 {
		t.Errorf("expected offset 0, got %d", offset)
	}
}

// TestSliceBoundsNilSafety tests that the analyzer doesn't crash on nil values
func TestSliceBoundsNilSafety(t *testing.T) {
	t.Run("extractBinOpBound with nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("extractBinOpBound panicked on nil input: %v", r)
			}
		}()
		_, _, _ = extractBinOpBound(nil)
	})

	t.Run("extractLenBound with nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("extractLenBound panicked on nil input: %v", r)
			}
		}()
		_, _, _ = extractLenBound(nil)
	})

	t.Run("extractBinOpBound with binop having nil X and Y", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("extractBinOpBound panicked on binop with nil X/Y: %v", r)
			}
		}()
		binop := &ssa.BinOp{Op: token.LSS}
		// X and Y are nil by default
		_, _, _ = extractBinOpBound(binop)
	})
}

// TestInvBound tests the invBound function
func TestInvBound(t *testing.T) {
	tests := []struct {
		name     string
		input    bound
		expected bound
	}{
		{"lowerUnbounded", lowerUnbounded, upperUnbounded},
		{"upperUnbounded", upperUnbounded, lowerUnbounded},
		{"upperBounded", upperBounded, unbounded},
		{"unbounded", unbounded, upperBounded},
		{"bounded", bounded, bounded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := invBound(tt.input)
			if result != tt.expected {
				t.Errorf("invBound(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
