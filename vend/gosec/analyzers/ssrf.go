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

// SSRF returns a configuration for detecting Server-Side Request Forgery vulnerabilities.
func SSRF() taint.Config {
	return taint.Config{
		Sources: []taint.Source{
			{Package: "net/http", Name: "Request", Pointer: true},
			{Package: "os", Name: "Args"},
			{Package: "os", Name: "Getenv"},
			{Package: "bufio", Name: "Reader", Pointer: true},
			{Package: "bufio", Name: "Scanner", Pointer: true},
			{Package: "os", Name: "File", Pointer: true},
		},
		Sinks: []taint.Sink{
			{Package: "net/http", Method: "Get"},
			{Package: "net/http", Method: "Post"},
			{Package: "net/http", Method: "Head"},
			{Package: "net/http", Method: "PostForm"},
			{Package: "net/http", Method: "NewRequest"},
			{Package: "net/http", Receiver: "Client", Method: "Do", Pointer: true},
			{Package: "net/http", Receiver: "Client", Method: "Get", Pointer: true},
			{Package: "net/http", Receiver: "Client", Method: "Post", Pointer: true},
			{Package: "net/http", Receiver: "Client", Method: "Head", Pointer: true},
			{Package: "net", Method: "Dial"},
			{Package: "net", Method: "DialTimeout"},
			{Package: "net", Method: "LookupHost"},
			{Package: "net/http/httputil", Method: "NewSingleHostReverseProxy"},
			{Package: "net/http/httputil", Receiver: "ReverseProxy", Method: "ServeHTTP", Pointer: true},
		},
	}
}

// newSSRFAnalyzer creates an analyzer for detecting SSRF vulnerabilities
// via taint analysis (G704)
func newSSRFAnalyzer(id string, description string) *analysis.Analyzer {
	config := SSRF()
	rule := SSRFRule
	rule.ID = id
	rule.Description = description
	return taint.NewGosecAnalyzer(&rule, &config)
}
