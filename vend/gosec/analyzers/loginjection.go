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

// LogInjection returns a configuration for detecting log injection vulnerabilities.
func LogInjection() taint.Config {
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
			{Package: "log", Method: "Print"},
			{Package: "log", Method: "Printf"},
			{Package: "log", Method: "Println"},
			{Package: "log", Method: "Fatal"},
			{Package: "log", Method: "Fatalf"},
			{Package: "log", Method: "Fatalln"},
			{Package: "log", Method: "Panic"},
			{Package: "log", Method: "Panicf"},
			{Package: "log", Method: "Panicln"},
			{Package: "log/slog", Method: "Info"},
			{Package: "log/slog", Method: "Warn"},
			{Package: "log/slog", Method: "Error"},
			{Package: "log/slog", Method: "Debug"},
		},
	}
}

// newLogInjectionAnalyzer creates an analyzer for detecting log injection vulnerabilities
// via taint analysis (G706)
func newLogInjectionAnalyzer(id string, description string) *analysis.Analyzer {
	config := LogInjection()
	rule := LogInjectionRule
	rule.ID = id
	rule.Description = description
	return taint.NewGosecAnalyzer(&rule, &config)
}
