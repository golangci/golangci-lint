package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG203 - Template checks
var SampleCodeG203 = []CodeSample{
	{[]string{`
// We assume that hardcoded template strings are safe as the programmer would
// need to be explicitly shooting themselves in the foot (as below)
package main

import (
	"html/template"
	"os"
)

const tmpl = ""

func main() {
	t := template.Must(template.New("ex").Parse(tmpl))
	v := map[string]interface{}{
		"Title":    "Test <b>World</b>",
		"Body":     template.HTML("<script>alert(1)</script>"),
	}
	t.Execute(os.Stdout, v)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
// Using a variable to initialize could potentially be dangerous. Under the
// current model this will likely produce some false positives.
package main

import (
	"html/template"
	"os"
)

const tmpl = ""

func main() {
	a := "something from another place"
	t := template.Must(template.New("ex").Parse(tmpl))
	v := map[string]interface{}{
		"Title":    "Test <b>World</b>",
		"Body":     template.HTML(a),
	}
	t.Execute(os.Stdout, v)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"html/template"
	"os"
)

const tmpl = ""

func main() {
	a := "something from another place"
	t := template.Must(template.New("ex").Parse(tmpl))
	v := map[string]interface{}{
		"Title":    "Test <b>World</b>",
		"Body":     template.JS(a),
	}
	t.Execute(os.Stdout, v)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"html/template"
	"os"
)

const tmpl = ""

func main() {
	a := "something from another place"
	t := template.Must(template.New("ex").Parse(tmpl))
	v := map[string]interface{}{
		"Title":    "Test <b>World</b>",
		"Body":     template.URL(a),
	}
	t.Execute(os.Stdout, v)
}
`}, 1, gosec.NewConfig()},
}
