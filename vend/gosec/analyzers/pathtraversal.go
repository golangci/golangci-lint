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
	"golang.org/x/tools/go/analysis"

	"github.com/securego/gosec/v2/taint"
)

// PathTraversal returns a configuration for detecting path traversal vulnerabilities.
func PathTraversal() taint.Config {
	return taint.Config{
		Sources: []taint.Source{
			{Package: "net/http", Name: "Request", Pointer: true},
			{Package: "net/url", Name: "URL", Pointer: true},
			{Package: "os", Name: "Args"},
			{Package: "os", Name: "Getenv"},
			{Package: "bufio", Name: "Reader", Pointer: true},
			{Package: "bufio", Name: "Scanner", Pointer: true},
			{Package: "os", Name: "File", Pointer: true},
		},
		Sinks: []taint.Sink{
			{Package: "os", Method: "Open"},
			{Package: "os", Method: "OpenFile"},
			{Package: "os", Method: "Create"},
			{Package: "os", Method: "ReadFile"},
			{Package: "os", Method: "WriteFile"},
			{Package: "os", Method: "Remove"},
			{Package: "os", Method: "RemoveAll"},
			{Package: "os", Method: "Rename"},
			{Package: "os", Method: "Mkdir"},
			{Package: "os", Method: "MkdirAll"},
			{Package: "os", Method: "Stat"},
			{Package: "os", Method: "Lstat"},
			{Package: "os", Method: "Chmod"},
			{Package: "os", Method: "Chown"},
			{Package: "io/ioutil", Method: "ReadFile"},
			{Package: "io/ioutil", Method: "WriteFile"},
			{Package: "io/ioutil", Method: "ReadDir"},
			{Package: "path/filepath", Method: "Walk"},
			{Package: "path/filepath", Method: "WalkDir"},
		},
	}
}

// newPathTraversalAnalyzer creates an analyzer for detecting path traversal vulnerabilities
// via taint analysis (G703)
func newPathTraversalAnalyzer(id string, description string) *analysis.Analyzer {
	config := PathTraversal()
	rule := PathTraversalRule
	rule.ID = id
	rule.Description = description
	return taint.NewGosecAnalyzer(&rule, &config)
}
