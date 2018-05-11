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
	"regexp"

	"github.com/GoASTScanner/gas"
)

type badTempFile struct {
	gas.MetaData
	calls gas.CallList
	args  *regexp.Regexp
}

func (t *badTempFile) ID() string {
	return t.MetaData.ID
}

func (t *badTempFile) Match(n ast.Node, c *gas.Context) (gi *gas.Issue, err error) {
	if node := t.calls.ContainsCallExpr(n, c); node != nil {
		if arg, e := gas.GetString(node.Args[0]); t.args.MatchString(arg) && e == nil {
			return gas.NewIssue(c, n, t.ID(), t.What, t.Severity, t.Confidence), nil
		}
	}
	return nil, nil
}

// NewBadTempFile detects direct writes to predictable path in temporary directory
func NewBadTempFile(id string, conf gas.Config) (gas.Rule, []ast.Node) {
	calls := gas.NewCallList()
	calls.Add("io/ioutil", "WriteFile")
	calls.Add("os", "Create")
	return &badTempFile{
		calls: calls,
		args:  regexp.MustCompile(`^/tmp/.*$|^/var/tmp/.*$`),
		MetaData: gas.MetaData{
			ID:         id,
			Severity:   gas.Medium,
			Confidence: gas.High,
			What:       "File creation in shared tmp directory without using ioutil.Tempfile",
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
