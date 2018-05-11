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
	"fmt"
	"go/ast"

	"github.com/GoASTScanner/gas"
)

type weakKeyStrength struct {
	gas.MetaData
	calls gas.CallList
	bits  int
}

func (w *weakKeyStrength) ID() string {
	return w.MetaData.ID
}

func (w *weakKeyStrength) Match(n ast.Node, c *gas.Context) (*gas.Issue, error) {
	if callExpr := w.calls.ContainsCallExpr(n, c); callExpr != nil {
		if bits, err := gas.GetInt(callExpr.Args[1]); err == nil && bits < (int64)(w.bits) {
			return gas.NewIssue(c, n, w.ID(), w.What, w.Severity, w.Confidence), nil
		}
	}
	return nil, nil
}

// NewWeakKeyStrength builds a rule that detects RSA keys < 2048 bits
func NewWeakKeyStrength(id string, conf gas.Config) (gas.Rule, []ast.Node) {
	calls := gas.NewCallList()
	calls.Add("crypto/rsa", "GenerateKey")
	bits := 2048
	return &weakKeyStrength{
		calls: calls,
		bits:  bits,
		MetaData: gas.MetaData{
			ID:         id,
			Severity:   gas.Medium,
			Confidence: gas.High,
			What:       fmt.Sprintf("RSA keys should be at least %d bits", bits),
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
