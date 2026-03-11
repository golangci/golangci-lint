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
	"fmt"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"

	"github.com/securego/gosec/v2/internal/ssautil"
	"github.com/securego/gosec/v2/issue"
)

// newConversionOverflowAnalyzer creates a new analysis.Analyzer for detecting integer overflows in conversions.
func newConversionOverflowAnalyzer(id string, description string) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     id,
		Doc:      description,
		Run:      runConversionOverflow,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	}
}

type conversionPair struct {
	src types.BasicKind
	dst types.BasicKind
}

type overflowState struct {
	*BaseAnalyzerState
	msgCache map[conversionPair]string
}

func newOverflowState(pass *analysis.Pass) *overflowState {
	return &overflowState{
		BaseAnalyzerState: NewBaseState(pass),
		msgCache:          make(map[conversionPair]string),
	}
}

// runConversionOverflow analyzes the SSA representation of the code to find potential integer overflows in type conversions.
func runConversionOverflow(pass *analysis.Pass) (any, error) {
	ssaResult, err := ssautil.GetSSAResult(pass)
	if err != nil {
		return nil, fmt.Errorf("building ssa representation: %w", err)
	}

	state := newOverflowState(pass)
	defer state.Release()
	issues := []*issue.Issue{}
	for _, mcall := range ssaResult.SSA.SrcFuncs {
		state.Reset()
		for _, block := range mcall.DomPreorder() {
			for _, instr := range block.Instrs {
				switch instr := instr.(type) {
				case *ssa.Convert:
					srcInfo, err := GetIntTypeInfo(instr.X.Type())
					if err != nil {
						continue
					}
					dstInfo, err := GetIntTypeInfo(instr.Type())
					if err != nil {
						continue
					}

					if hasOverflow(srcInfo, dstInfo) {
						if state.isSafeConversion(instr, dstInfo) {
							continue
						}

						srcBasic, _ := instr.X.Type().Underlying().(*types.Basic)
						dstBasic, _ := instr.Type().Underlying().(*types.Basic)

						if srcBasic == nil || dstBasic == nil {
							continue
						}

						pair := conversionPair{
							src: srcBasic.Kind(),
							dst: dstBasic.Kind(),
						}
						msg, ok := state.msgCache[pair]
						if !ok {
							msg = fmt.Sprintf("integer overflow conversion %s -> %s", srcBasic.Name(), dstBasic.Name())
							state.msgCache[pair] = msg
						}

						issues = append(issues, newIssue(pass.Analyzer.Name,
							msg,
							pass.Fset,
							instr.Pos(),
							issue.High,
							issue.Medium,
						))
					}
				}
			}
		}
	}

	if len(issues) > 0 {
		return issues, nil
	}
	return nil, nil
}

// isSafeConversion checks if a specific conversion instruction is safe from overflow, considering logic and constraints.
func (s *overflowState) isSafeConversion(instr *ssa.Convert, dstInt IntTypeInfo) bool {
	// Check for constant conversions.
	if constVal, ok := instr.X.(*ssa.Const); ok {
		if IsConstantInTypeRange(constVal, dstInt) {
			return true
		}
	}

	// Check for explicit range checks.
	if s.hasRangeCheck(instr.X, dstInt, instr.Block()) {
		return true
	}
	return false
}

func hasOverflow(srcInfo, dstInfo IntTypeInfo) bool {
	return srcInfo.Min < dstInfo.Min || srcInfo.Max > dstInfo.Max
}

// hasRangeCheck determines if there is a valid range check for the given value that ensures safety.
func (s *overflowState) hasRangeCheck(v ssa.Value, dstInt IntTypeInfo, block *ssa.BasicBlock) bool {
	// Clear visited map for new resolution
	clear(s.Visited)

	res := s.Analyzer.ResolveRange(v, block)
	defer s.Analyzer.releaseResult(res)

	// Check for explicit values
	if ExplicitValsInRange(res.explicitPositiveVals, res.explicitNegativeVals, dstInt) {
		return true
	}

	// Check all predecessors for OR support.
	if len(block.Preds) > 1 {
		allPredsSafe := true
		for _, pred := range block.Preds {
			if !s.isSafeFromPredecessor(v, dstInt, pred, block) {
				allPredsSafe = false
				break
			}
		}
		if allPredsSafe {
			return true
		}
	}

	// Relax requirement: If we have a definitive range (both set) and it's safe,
	// we allow it even if not explicitly "checked" by an IF,
	// because definition-based ranges (like constants or arithmetic on constants) are certain.
	isDefinitiveSafe := res.minValueSet && res.maxValueSet

	if !res.isRangeCheck && !isDefinitiveSafe {
		return false
	}

	return s.validateRangeLimits(v, res, dstInt)
}

func (s *overflowState) validateRangeLimits(v ssa.Value, res *rangeResult, dstInt IntTypeInfo) bool {
	minValue, minValueSet, maxValue, maxValueSet := res.minValue, res.minValueSet, res.maxValue, res.maxValueSet
	isSrcUnsigned := isUint(v)

	// Check for impossible ranges (disjoint)
	if !isSrcUnsigned {
		if minValueSet && maxValueSet && toInt64(minValue) > toInt64(maxValue) {
			return true
		}
	}
	if isSrcUnsigned && minValueSet && maxValueSet && minValue > maxValue {
		return true
	}

	srcInt, err := GetIntTypeInfo(v.Type())
	if err != nil {
		return false
	}

	if dstInt.Signed {
		if isSrcUnsigned {
			return maxValueSet && maxValue <= dstInt.Max
		}
		minSafe := true
		if srcInt.Min < dstInt.Min {
			minSafe = minValueSet && toInt64(minValue) >= dstInt.Min
		}
		maxSafe := true
		if srcInt.Max > dstInt.Max {
			maxSafe = maxValueSet && toInt64(maxValue) <= toInt64(dstInt.Max)
		}
		return minSafe && maxSafe
	}
	if isSrcUnsigned {
		return maxValueSet && maxValue <= dstInt.Max
	}
	minSafe := true
	if srcInt.Min < 0 {
		minSafe = minValueSet && toInt64(minValue) >= 0
	}
	maxSafe := true
	if srcInt.Max > dstInt.Max {
		maxSafe = maxValueSet && maxValue <= dstInt.Max
	}
	return minSafe && maxSafe
}

func (s *overflowState) isSafeFromPredecessor(v ssa.Value, dstInt IntTypeInfo, pred *ssa.BasicBlock, targetBlock *ssa.BasicBlock) bool {
	if vIf, ok := pred.Instrs[len(pred.Instrs)-1].(*ssa.If); ok {
		isSrcUnsigned := isUint(v)
		for i, succ := range pred.Succs {
			if succ == targetBlock {
				// We took this specific edge.
				result := s.Analyzer.getResultRangeForIfEdge(vIf, i == 0, v)
				defer s.Analyzer.releaseResult(result)

				if result.isRangeCheck {
					var safe bool
					if dstInt.Signed {
						if isSrcUnsigned {
							safe = result.maxValueSet && result.maxValue <= dstInt.Max
						} else {
							safe = (result.minValueSet && toInt64(result.minValue) >= dstInt.Min) && (result.maxValueSet && toInt64(result.maxValue) <= toInt64(dstInt.Max))
						}
					} else {
						if isSrcUnsigned {
							safe = result.maxValueSet && result.maxValue <= dstInt.Max
						} else {
							safe = (result.minValueSet && toInt64(result.minValue) >= 0) && (result.maxValueSet && result.maxValue <= dstInt.Max)
						}
					}
					if safe {
						return true
					}
				}
			}
		}
	}
	return false
}
