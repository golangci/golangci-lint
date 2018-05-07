// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
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

package rules

import (
	"go/ast"
	"go/types"

	"github.com/GoASTScanner/gas"
)

type subprocess struct {
	gas.MetaData
	gas.CallList
}

func (r *subprocess) ID() string {
	return r.MetaData.ID
}

// TODO(gm) The only real potential for command injection with a Go project
// is something like this:
//
// syscall.Exec("/bin/sh", []string{"-c", tainted})
//
// E.g. Input is correctly escaped but the execution context being used
// is unsafe. For example:
//
// syscall.Exec("echo", "foobar" + tainted)
func (r *subprocess) Match(n ast.Node, c *gas.Context) (*gas.Issue, error) {
	if node := r.ContainsCallExpr(n, c); node != nil {
		for _, arg := range node.Args {
			if ident, ok := arg.(*ast.Ident); ok {
				obj := c.Info.ObjectOf(ident)
				if _, ok := obj.(*types.Var); ok && !gas.TryResolve(ident, c) {
					return gas.NewIssue(c, n, r.ID(), "Subprocess launched with variable", gas.Medium, gas.High), nil
				}
			}
		}
		return gas.NewIssue(c, n, r.ID(), "Subprocess launching should be audited", gas.Low, gas.High), nil
	}
	return nil, nil
}

// NewSubproc detects cases where we are forking out to an external process
func NewSubproc(id string, conf gas.Config) (gas.Rule, []ast.Node) {
	rule := &subprocess{gas.MetaData{ID: id}, gas.NewCallList()}
	rule.Add("os/exec", "Command")
	rule.Add("syscall", "Exec")
	return rule, []ast.Node{(*ast.CallExpr)(nil)}
}
