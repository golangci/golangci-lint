package lintpack

import (
	"go/ast"
)

// CheckerInfo holds checker metadata and structured documentation.
type CheckerInfo struct {
	// Name is a checker name.
	Name string

	// Tags is a list of labels that can be used to enable or disable checker.
	// Common tags are "experimental" and "performance".
	Tags []string

	// Summary is a short one sentence description.
	// Should not end with a period.
	Summary string

	// Details extends summary with additional info. Optional.
	Details string

	// Before is a code snippet of code that will violate rule.
	Before string

	// After is a code snippet of fixed code that complies to the rule.
	After string

	// Note is an optional caution message or advice.
	Note string
}

// Warning represents issue that is found by checker.
type Warning struct {
	// Node is an AST node that caused warning to trigger.
	// Can be used to obtain proper error location.
	Node ast.Node

	// Text is warning message without source location info.
	Text string
}
