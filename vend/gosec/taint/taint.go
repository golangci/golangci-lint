// Package taint provides a minimal taint analysis engine for gosec.
// It tracks data flow from sources (user input) to sinks (dangerous functions)
// using SSA form and call graph analysis.
//
// This implementation uses only golang.org/x/tools packages which gosec
// already depends on - no external dependencies required.
//
// Inspired by:
//   - github.com/google/capslock (call graph traversal pattern)
//   - gosec issue #1160 (requirements)
package taint

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/ssa"
)

// maxTaintDepth limits recursion depth to prevent stack overflow on large codebases
const maxTaintDepth = 50

// Source defines where tainted data originates.
// Format: "package/path.TypeOrFunc" or "*package/path.Type" for pointer types.
type Source struct {
	// Package is the import path of the package containing the source (e.g., "net/http")
	Package string
	// Name is the type or function name that produces tainted data (e.g., "Request" for type, "Get" for function)
	Name string
	// Pointer indicates whether the source is a pointer type (true for *Type)
	Pointer bool
}

// Sink defines a dangerous function that should not receive tainted data.
// Format: "(*package/path.Type).Method" or "package/path.Func"
type Sink struct {
	// Package is the import path of the package containing the sink (e.g., "database/sql")
	Package string
	// Receiver is the type name for methods (e.g., "DB"), or empty for package-level functions
	Receiver string
	// Method is the function or method name that represents the sink (e.g., "Query")
	Method string
	// Pointer indicates whether the receiver is a pointer type (true for *Type methods)
	Pointer bool
	// CheckArgs specifies which argument positions to check for taint (0-indexed).
	// For method calls, Args[0] is the receiver.
	// If nil or empty, all arguments are checked.
	// Examples:
	//   - SQL methods: [1] - only check query string (Args[1]), skip receiver
	//   - fmt.Fprintf: [1,2,3,...] - skip writer (Args[0]), check format and data
	CheckArgs []int
}

// Result represents a detected taint flow from source to sink.
type Result struct {
	// Source is the origin of the tainted data
	Source Source
	// Sink is the dangerous function that receives the tainted data
	Sink Sink
	// SinkPos is the source code position of the sink call
	SinkPos token.Pos
	// Path is the sequence of functions from entry point to the sink
	Path []*ssa.Function
}

// Config holds taint analysis configuration.
type Config struct {
	// Sources is the list of data origins that produce tainted values
	Sources []Source
	// Sinks is the list of dangerous functions that should not receive tainted data
	Sinks []Sink
}

// Analyzer performs taint analysis on SSA programs.
type Analyzer struct {
	config    *Config
	sources   map[string]Source // keyed by full type string
	sinks     map[string]Sink   // keyed by full function string
	callGraph *callgraph.Graph
}

// New creates a new taint analyzer with the given configuration.
func New(config *Config) *Analyzer {
	a := &Analyzer{
		config:  config,
		sources: make(map[string]Source),
		sinks:   make(map[string]Sink),
	}

	// Index sources for fast lookup
	for _, src := range config.Sources {
		key := formatSourceKey(src)
		a.sources[key] = src
	}

	// Index sinks for fast lookup
	for _, sink := range config.Sinks {
		key := formatSinkKey(sink)
		a.sinks[key] = sink
	}

	return a
}

// formatSourceKey creates a lookup key for a source.
func formatSourceKey(src Source) string {
	key := src.Package + "." + src.Name
	if src.Pointer {
		key = "*" + key
	}
	return key
}

// formatSinkKey creates a lookup key for a sink.
func formatSinkKey(sink Sink) string {
	if sink.Receiver == "" {
		return sink.Package + "." + sink.Method
	}
	recv := sink.Package + "." + sink.Receiver
	if sink.Pointer {
		recv = "*" + recv
	}
	return "(" + recv + ")." + sink.Method
}

// Analyze performs taint analysis on the given SSA program.
// It returns all detected taint flows from sources to sinks.
func (a *Analyzer) Analyze(prog *ssa.Program, srcFuncs []*ssa.Function) []Result {
	if len(srcFuncs) == 0 {
		return nil
	}

	// Build call graph using Class Hierarchy Analysis (CHA).
	// CHA is fast and sound (no false negatives) but may have false positives.
	// For more precision, use VTA (Variable Type Analysis) instead.
	a.callGraph = cha.CallGraph(prog)

	var results []Result

	// Find all sink calls in the program
	for _, fn := range srcFuncs {
		results = append(results, a.analyzeFunctionSinks(fn)...)
	}

	return results
}

// analyzeFunctionSinks finds sink calls in a function and traces taint.
func (a *Analyzer) analyzeFunctionSinks(fn *ssa.Function) []Result {
	if fn == nil || fn.Blocks == nil {
		return nil
	}

	var results []Result

	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			call, ok := instr.(*ssa.Call)
			if !ok {
				continue
			}

			// Check if this call is a sink
			sink, isSink := a.isSinkCall(call)
			if !isSink {
				continue
			}

			// Determine which arguments to check for taint
			var argsToCheck []ssa.Value

			if len(sink.CheckArgs) > 0 {
				// Sink specifies which argument positions to check
				for _, idx := range sink.CheckArgs {
					if idx < len(call.Call.Args) {
						argsToCheck = append(argsToCheck, call.Call.Args[idx])
					}
				}
			} else {
				// No CheckArgs specified: check all arguments
				argsToCheck = call.Call.Args
			}

			// Check if any of the specified arguments are tainted
			for _, arg := range argsToCheck {
				if a.isTainted(arg, fn, make(map[ssa.Value]bool), 0) {
					results = append(results, Result{
						Sink:    sink,
						SinkPos: call.Pos(),
						Path:    a.buildPath(fn),
					})
					break
				}
			}
		}
	}

	return results
}

// isSinkCall checks if a call instruction is a sink and returns the sink info.
func (a *Analyzer) isSinkCall(call *ssa.Call) (Sink, bool) {
	// Try to get receiver info first (works for both concrete and interface calls)
	var pkg, receiverName, methodName string
	var isPointer bool

	// Check for method call (invoke or static with receiver)
	if call.Call.IsInvoke() {
		// Interface method call - receiver is in Call.Value, not Args
		if call.Call.Value != nil {
			recvType := call.Call.Value.Type()
			methodName = call.Call.Method.Name()

			// For interface calls, the type is usually a Named type pointing to the interface
			if named, ok := recvType.(*types.Named); ok {
				receiverName = named.Obj().Name()
				if pkgObj := named.Obj(); pkgObj != nil && pkgObj.Pkg() != nil {
					pkg = pkgObj.Pkg().Path()
				}
			}

			// Match against sinks (interface methods don't have Pointer field usually)
			for _, sink := range a.sinks {
				if sink.Package == pkg && sink.Receiver == receiverName && sink.Method == methodName {
					return sink, true
				}
			}
		}
	}

	// Try static callee (for non-interface method calls and functions)
	callee := call.Call.StaticCallee()
	if callee != nil {
		if callee.Pkg != nil && callee.Pkg.Pkg != nil {
			pkg = callee.Pkg.Pkg.Path()
		}
		methodName = callee.Name()

		// Check if it has a receiver (method call)
		if recv := callee.Signature.Recv(); recv != nil {
			recvType := recv.Type()
			if named, ok := recvType.(*types.Named); ok {
				receiverName = named.Obj().Name()
			}
			if ptr, ok := recvType.(*types.Pointer); ok {
				isPointer = true
				if named, ok := ptr.Elem().(*types.Named); ok {
					receiverName = named.Obj().Name()
				}
			}
		}
	}

	// Match against configured sinks
	for _, sink := range a.sinks {
		// Package must match
		if sink.Package != pkg {
			continue
		}

		// For method sinks (with receiver)
		if sink.Receiver != "" {
			if sink.Receiver == receiverName && sink.Method == methodName && sink.Pointer == isPointer {
				return sink, true
			}
		} else {
			// For function sinks (no receiver)
			if sink.Method == methodName && receiverName == "" {
				return sink, true
			}
		}
	}

	return Sink{}, false
}

// isTainted recursively checks if a value is tainted (originates from a source).
func (a *Analyzer) isTainted(v ssa.Value, fn *ssa.Function, visited map[ssa.Value]bool, depth int) bool {
	if v == nil {
		return false
	}

	// Prevent stack overflow on large codebases
	if depth > maxTaintDepth {
		return false
	}

	// Prevent infinite recursion
	if visited[v] {
		return false
	}
	visited[v] = true

	// Check if this value's type is a source
	if a.isSourceType(v.Type()) {
		return true
	}

	// Trace back through SSA instructions
	switch val := v.(type) {
	case *ssa.Parameter:
		// Parameters can be tainted if the function is called with tainted args
		return a.isParameterTainted(val, fn, visited, depth+1)

	case *ssa.Call:
		// Check if calling a method on a tainted receiver propagates taint
		if val.Call.IsInvoke() {
			// Method call on interface: check if receiver is tainted
			if len(val.Call.Args) > 0 && a.isTainted(val.Call.Args[0], fn, visited, depth+1) {
				return true
			}
		} else if callee := val.Call.StaticCallee(); callee != nil && callee.Signature.Recv() != nil {
			// Method call with receiver: check if receiver is tainted
			if len(val.Call.Args) > 0 && a.isTainted(val.Call.Args[0], fn, visited, depth+1) {
				return true
			}
		}

		// Check if the call returns a tainted type
		if a.isSourceCall(val) {
			return true
		}
		// Check if the receiver (for method calls) is tainted
		if val.Call.Value != nil {
			if a.isTainted(val.Call.Value, fn, visited, depth+1) {
				return true
			}
		}
		// Check if any argument to this call is tainted
		// This handles conversions like []byte(taintedString)
		for _, arg := range val.Call.Args {
			if a.isTainted(arg, fn, visited, depth+1) {
				return true
			}
		}
		// Check for builtin conversions (like string to []byte)
		if builtin, ok := val.Call.Value.(*ssa.Builtin); ok {
			_ = builtin // Builtins like "append", "copy", etc.
			// For builtins, if any argument is tainted, result is tainted
			for _, arg := range val.Call.Args {
				if a.isTainted(arg, fn, visited, depth+1) {
					return true
				}
			}
		}

	case *ssa.FieldAddr:
		// Field access on a tainted struct
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.IndexAddr:
		// Index into a tainted slice/array
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.UnOp:
		// Unary operation (like pointer dereference)
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.BinOp:
		// Binary operation - tainted if either operand is tainted
		return a.isTainted(val.X, fn, visited, depth+1) || a.isTainted(val.Y, fn, visited, depth+1)

	case *ssa.Phi:
		// Phi node - tainted if any edge is tainted
		for _, edge := range val.Edges {
			if a.isTainted(edge, fn, visited, depth+1) {
				return true
			}
		}

	case *ssa.Extract:
		// Extract from tuple - check the tuple
		return a.isTainted(val.Tuple, fn, visited, depth+1)

	case *ssa.TypeAssert:
		// Type assertion - check the underlying value
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.MakeInterface:
		// Interface creation - check the underlying value
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.Slice:
		// Slice operation - check the sliced value
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.Convert:
		// Type conversion - check the converted value
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.ChangeType:
		// Type change - check the underlying value
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.Alloc:
		// Allocation - check referrers for assignments
		for _, ref := range *val.Referrers() {
			// Direct stores to the allocation
			if store, ok := ref.(*ssa.Store); ok {
				if a.isTainted(store.Val, fn, visited, depth+1) {
					return true
				}
			}
			// For arrays/slices, check stores to indexed addresses (e.g., varargs)
			if indexAddr, ok := ref.(*ssa.IndexAddr); ok {
				if indexRefs := indexAddr.Referrers(); indexRefs != nil {
					for _, indexRef := range *indexRefs {
						if store, ok := indexRef.(*ssa.Store); ok {
							if a.isTainted(store.Val, fn, visited, depth+1) {
								return true
							}
						}
					}
				}
			}
		}

	case *ssa.Lookup:
		// Map/string lookup - check the map/string
		return a.isTainted(val.X, fn, visited, depth+1)

	case *ssa.MakeSlice:
		// MakeSlice - check if it's being populated with tainted data
		// This handles cases like []byte(taintedString)
		if refs := val.Referrers(); refs != nil {
			for _, ref := range *refs {
				// Check stores into the slice
				if store, ok := ref.(*ssa.Store); ok {
					if a.isTainted(store.Val, fn, visited, depth+1) {
						return true
					}
				}
				// Check if used in a call that populates it (like copy())
				if call, ok := ref.(*ssa.Call); ok {
					for _, arg := range call.Call.Args {
						if arg == val {
							continue // Skip the slice itself
						}
						if a.isTainted(arg, fn, visited, depth+1) {
							return true
						}
					}
				}
			}
		}
		return false

	case *ssa.MakeMap, *ssa.MakeChan:
		// New maps/channels are not tainted by default
		return false

	case *ssa.Const:
		// Constants are never tainted
		return false

	case *ssa.Global:
		// Global variables - check if they're a known source (e.g., os.Args)
		if a.isSourceType(val.Type()) {
			return true
		}
		// Check if the global variable itself is configured as a source
		if val.Pkg != nil && val.Pkg.Pkg != nil {
			globalKey := val.Pkg.Pkg.Path() + "." + val.Name()
			if _, ok := a.sources[globalKey]; ok {
				return true
			}
		}
		return false

	default:
		// Unhandled SSA instruction type - be conservative and don't propagate taint
		// to avoid false positives, but this might cause false negatives
		return false
	}

	return false
}

// isSourceType checks if a type matches any configured source.
func (a *Analyzer) isSourceType(t types.Type) bool {
	if t == nil {
		return false
	}

	typeStr := t.String()

	// Direct match
	if _, ok := a.sources[typeStr]; ok {
		return true
	}

	// Check underlying type for named types
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil {
			key := obj.Pkg().Path() + "." + obj.Name()
			if _, ok := a.sources[key]; ok {
				return true
			}
			// Check pointer variant
			if _, ok := a.sources["*"+key]; ok {
				return true
			}
		}
	}

	// Check pointer types
	if ptr, ok := t.(*types.Pointer); ok {
		return a.isSourceType(ptr.Elem())
	}

	return false
}

// isSourceCall checks if a call returns a value from a source function.
func (a *Analyzer) isSourceCall(call *ssa.Call) bool {
	callee := call.Call.StaticCallee()
	if callee == nil {
		return false
	}

	// Check if return type is a source
	if a.isSourceType(call.Type()) {
		return true
	}

	// Check if the function itself is a source (e.g., os.Getenv, os.ReadFile)
	if callee.Pkg != nil && callee.Pkg.Pkg != nil {
		pkg := callee.Pkg.Pkg.Path()
		funcKey := pkg + "." + callee.Name()
		if _, ok := a.sources[funcKey]; ok {
			return true
		}
	}

	return false
}

// isParameterTainted checks if a function parameter receives tainted data.
func (a *Analyzer) isParameterTainted(param *ssa.Parameter, fn *ssa.Function, visited map[ssa.Value]bool, depth int) bool {
	// Prevent stack overflow
	if depth > maxTaintDepth {
		return false
	}

	// Check if parameter type is a source
	if a.isSourceType(param.Type()) {
		return true
	}

	// Use call graph to find callers and check their arguments
	if a.callGraph == nil {
		return false
	}

	node := a.callGraph.Nodes[fn]
	if node == nil {
		return false
	}

	paramIdx := -1
	for i, p := range fn.Params {
		if p == param {
			paramIdx = i
			break
		}
	}

	if paramIdx < 0 {
		return false
	}

	// Check each caller
	for _, inEdge := range node.In {
		site := inEdge.Site
		if site == nil {
			continue
		}

		callArgs := site.Common().Args

		// Adjust for receiver
		if fn.Signature.Recv() != nil {
			paramIdx++
		}

		if paramIdx < len(callArgs) {
			if a.isTainted(callArgs[paramIdx], inEdge.Caller.Func, visited, depth+1) {
				return true
			}
		}
	}

	return false
}

// buildPath constructs the call path from entry point to the sink.
func (a *Analyzer) buildPath(fn *ssa.Function) []*ssa.Function {
	if a.callGraph == nil {
		return []*ssa.Function{fn}
	}

	// BFS to find path from root to this function
	path := []*ssa.Function{fn}

	node := a.callGraph.Nodes[fn]
	if node == nil {
		return path
	}

	// Simple path: just trace callers up
	visited := make(map[*ssa.Function]bool)
	current := node

	for current != nil && len(current.In) > 0 {
		if visited[current.Func] {
			break
		}
		visited[current.Func] = true

		caller := current.In[0].Caller
		if caller == nil || caller.Func == nil {
			break
		}

		path = append([]*ssa.Function{caller.Func}, path...)
		current = caller
	}

	return path
}
