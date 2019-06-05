[![Build Status](https://travis-ci.org/valyala/quicktemplate.svg)](https://travis-ci.org/valyala/quicktemplate)
[![GoDoc](https://godoc.org/github.com/valyala/quicktemplate?status.svg)](http://godoc.org/github.com/valyala/quicktemplate)
[![Go Report Card](https://goreportcard.com/badge/github.com/valyala/quicktemplate)](https://goreportcard.com/report/github.com/valyala/quicktemplate)

# quicktemplate

A fast, powerful, yet easy to use template engine for Go.
Inspired by the [Mako templates](http://www.makotemplates.org/) philosophy.

# Features

  * [Extremely fast](#performance-comparison-with-htmltemplate).
    Templates are converted into Go code and then compiled.
  * Quicktemplate syntax is very close to Go - there is no need to learn
    yet another template language before starting to use quicktemplate.
  * Almost all bugs are caught during template compilation, so production
    suffers less from template-related bugs.
  * Easy to use. See [quickstart](#quick-start) and [examples](https://github.com/valyala/quicktemplate/tree/master/examples)
    for details.
  * Powerful. Arbitrary Go code may be embedded into and mixed with templates.
    Be careful with this power - do not query the database and/or external resources from
    templates unless you miss the PHP way in Go :) This power is mostly for
    arbitrary data transformations.
  * Easy to use template inheritance powered by [Go interfaces](https://golang.org/doc/effective_go.html#interfaces).
    See [this example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver) for details.
  * Templates are compiled into a single binary, so there is no need to copy
    template files to the server.

# Drawbacks

  * Templates cannot be updated on the fly on the server, since they
    are compiled into a single binary.
    Take a look at [fasttemplate](https://github.com/valyala/fasttemplate)
    if you need a fast template engine for simple dynamically updated templates.

# Performance comparison with html/template

Quicktemplate is more than 20x faster than [html/template](https://golang.org/pkg/html/template/).
The following simple template is used in the benchmark:

  * [html/template version](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.tpl)
  * [quicktemplate version](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.qtpl)

Benchmark results:

```
$ go test -bench='Benchmark(Quick|HTML)Template' -benchmem github.com/valyala/quicktemplate/tests
BenchmarkQuickTemplate1-4                 	10000000	       120 ns/op	       0 B/op	       0 allocs/op
BenchmarkQuickTemplate10-4                	 3000000	       441 ns/op	       0 B/op	       0 allocs/op
BenchmarkQuickTemplate100-4               	  300000	      3945 ns/op	       0 B/op	       0 allocs/op
BenchmarkHTMLTemplate1-4                  	  500000	      2501 ns/op	     752 B/op	      23 allocs/op
BenchmarkHTMLTemplate10-4                 	  100000	     12442 ns/op	    3521 B/op	     117 allocs/op
BenchmarkHTMLTemplate100-4                	   10000	    123392 ns/op	   34498 B/op	    1152 allocs/op
```

[goTemplateBenchmark](https://github.com/SlinSo/goTemplateBenchmark) compares QuickTemplate with numerous Go templating packages. QuickTemplate performs favorably.

# Security

  * All template placeholders are HTML-escaped by default.
  * Template placeholders for JSON strings prevent from `</script>`-based
    XSS attacks:

  ```qtpl
  {% func FailedXSS() %}
  <script>
      var s = {%q= "</script><script>alert('you pwned!')" %};
  </script>
  {% endfunc %}
  ```

# Examples

See [examples](https://github.com/valyala/quicktemplate/tree/master/examples).

# Quick start

First of all, install the `quicktemplate` package
and [quicktemplate compiler](https://github.com/valyala/quicktemplate/tree/master/qtc) (`qtc`):

```
go get -u github.com/valyala/quicktemplate
go get -u github.com/valyala/quicktemplate/qtc
```

Let's start with a minimal template example:

```qtpl
All text outside function templates is treated as comments,
i.e. it is just ignored by quicktemplate compiler (`qtc`). It is for humans.

Hello is a simple template function.
{% func Hello(name string) %}
	Hello, {%s name %}!
{% endfunc %}
```

Save this file into a `templates` folder under the name `hello.qtpl`
and run `qtc` inside this folder.

If everything went OK, `hello.qtpl.go` file should appear in the `templates` folder.
This file contains Go code for `hello.qtpl`. Let's use it!

Create a file main.go outside `templates` folder and put the following
code there:

```go
package main

import (
	"fmt"

	"./templates"
)

func main() {
	fmt.Printf("%s\n", templates.Hello("Foo"))
	fmt.Printf("%s\n", templates.Hello("Bar"))
}
```

Then issue `go run`. If everything went OK, you'll see something like this:

```

	Hello, Foo!


	Hello, Bar!

```

Let's create more a complex template which calls other template functions,
contains loops, conditions, breaks, continues and returns.
Put the following template into `templates/greetings.qtpl`:

```qtpl

Greetings greets up to 42 names.
It also greets John differently comparing to others.
{% func Greetings(names []string) %}
	{% if len(names) == 0 %}
		Nobody to greet :(
		{% return %}
	{% endif %}

	{% for i, name := range names %}
		{% if i == 42 %}
			I'm tired to greet so many people...
			{% break %}
		{% elseif name == "John" %}
			{%= sayHi("Mr. " + name) %}
			{% continue %}
		{% else %}
			{%= Hello(name) %}
		{% endif %}
	{% endfor %}
{% endfunc %}

sayHi is unexported, since it starts with lowercase letter.
{% func sayHi(name string) %}
	Hi, {%s name %}
{% endfunc %}

Note that every template file may contain an arbitrary number
of template functions. For instance, this file contains Greetings and sayHi
functions.
```

Run `qtc` inside `templates` folder. Now the folder should contain
two files with Go code: `hello.qtpl.go` and `greetings.qtpl.go`. These files
form a single `templates` Go package. Template functions and other template
stuff is shared between template files located in the same folder.
So `Hello` template function may be used inside `greetings.qtpl` while
it is defined in `hello.qtpl`.
Moreover, the folder may contain ordinary Go files, so its contents may
be used inside templates and vice versa.
The package name inside template files may be overriden
with `{% package packageName %}`.

Now put the following code into `main.go`:

```go
package main

import (
	"bytes"
	"fmt"

	"./templates"
)

func main() {
	names := []string{"Kate", "Go", "John", "Brad"}

	// qtc creates Write* function for each template function.
	// Such functions accept io.Writer as first parameter:
	var buf bytes.Buffer
	templates.WriteGreetings(&buf, names)

	fmt.Printf("buf=\n%s", buf.Bytes())
}
```

Careful readers may notice different output tags were used in these
templates: `{%s name %}` and `{%= Hello(name) %}`. What's the difference?
The `{%s x %}` is used for printing HTML-safe strings, while `{%= F() %}`
is used for embedding template function calls. Quicktemplate supports also
other output tags:

  * `{%d num %}` for integers.
  * `{%f float %}` for float64.
    Floating point precision may be set via `{%f.precision float %}`.
    For example, `{%f.2 1.2345 %}` outputs `1.23`.
  * `{%z bytes %}` for byte slices.
  * `{%q str %}` and `{%qz bytes %}` for JSON-compatible quoted strings.
  * `{%j str %}` and `{%jz bytes %}` for embedding str into a JSON string. Unlike `{%q str %}`,
    it doesn't quote the string.
  * `{%u str %}` and `{%uz bytes %}` for [URL encoding](https://en.wikipedia.org/wiki/Percent-encoding)
    the given str.
  * `{%v anything %}` is equivalent to `%v` in [printf-like functions](https://golang.org/pkg/fmt/).

All the output tags except `{%= F() %}` produce HTML-safe output, i.e. they
escape `<` to `&lt;`, `>` to `&gt;`, etc. If you don't want HTML-safe output,
then just put `=` after the tag. For example: `{%s= "<h1>This h1 won't be escaped</h1>" %}`.

As you may notice `{%= F() %}` and `{%s= F() %}` produce the same output for `{% func F() %}`.
But the first one is optimized for speed - it avoids memory allocations and copies.
It is therefore recommended to stick to it when embedding template function calls.

Additionally, the following extensions are supported for `{%= F() %}`:

  * `{%=h F() %}` produces html-escaped output.
  * `{%=u F() %}` produces [URL-encoded](https://en.wikipedia.org/wiki/Percent-encoding) output.
  * `{%=q F() %}` produces quoted json string.
  * `{%=j F() %}` produces json string without quotes.
  * `{%=uh F() %}` produces html-safe URL-encoded output.
  * `{%=qh F() %}` produces html-safe quoted json string.
  * `{%=jh F() %}` produces html-safe json string without quotes.

All output tags except `{%= F() %}` family may contain arbitrary valid
Go expressions instead of just an identifier. For example:

```qtpl
Import fmt for fmt.Sprintf()
{% import "fmt" %}

FmtFunc uses fmt.Sprintf() inside output tag
{% func FmtFunc(s string) %}
	{%s fmt.Sprintf("FmtFunc accepted %q string", s) %}
{% endfunc %}
```

There are other useful tags supported by quicktemplate:

  * `{% comment %}`

    ```qtpl
    {% comment %}
        This is a comment. It won't trap into the output.
        It may contain {% arbitrary tags %}. They are just ignored.
    {% endcomment %}
    ```

  * `{% plain %}`

    ```qtpl
    {% plain %}
        Tags will {% trap into %} the output {% unmodified %}.
        Plain block may contain invalid and {% incomplete tags.
    {% endplain %}
    ```

  * `{% collapsespace %}`

    ```qtpl
    {% collapsespace %}
        <div>
            <div>space between lines</div>
               and {%s "tags" %}
             <div>is collapsed into a single space
             unless{% newline %}or{% space %}is used</div>
        </div>
    {% endcollapsespace %}
    ```

    Is converted into:

    ```
    <div> <div>space between lines</div> and tags <div>is collapsed into a single space unless
    or is used</div> </div>
    ```

  * `{% stripspace %}`

    ```qtpl
    {% stripspace %}
         <div>
             <div>space between lines</div>
                and {%s " tags" %}
             <div>is removed unless{% newline %}or{% space %}is used</div>
         </div>
    {% endstripspace %}
    ```

    Is converted into:

    ```
    <div><div>space between lines</div>and tags<div>is removed unless
    or is used</div></div>
    ```

  * `{% switch %}`, `{% case %}` and `{% default %}`:


    ```qtpl
    1 + 1 =
    {% switch 1+1 %}
    {% case 2 %}
	2?
    {% case 42 %}
	42!
    {% default %}
        I don't know :(
    {% endswitch %}
    ```

  * `{% code %}`:

    ```qtpl
    {% code
    // arbitrary Go code may be embedded here!
    type FooArg struct {
        Name string
        Age int
    }
    %}
    ```

  * `{% package %}`:

    ```qtpl
    Override default package name with the custom name
    {% package customPackageName %}
    ```

  * `{% import %}`:

    ```qtpl
    Import external packages.
    {% import "foo/bar" %}
    {% import (
        "foo"
        bar "baz/baa"
    ) %}
    ```

  * `{% cat "/path/to/file" %}`:

    ```qtpl
    Cat emits the given file contents as a plaintext:
    {% func passwords() %}
        /etc/passwd contents:
        {% cat "/etc/passwd" %}
    {% endfunc %}
    ```

  * `{% interface %}`:

    ```qtpl
    Interfaces allow powerful templates' inheritance
    {%
    interface Page {
        Title()
        Body(s string, n int)
        Footer()
    }
    %}

    PrintPage prints Page
    {% func PrintPage(p Page) %}
        <html>
            <head><title>{%= p.Title() %}</title></head>
            <body>
                <div>{%= p.Body("foo", 42) %}</div>
                <div>{%= p.Footer() %}</div>
            </body>
        </html>
    {% endfunc %}

    Base page implementation
    {% code
    type BasePage struct {
        TitleStr string
        FooterStr string
    }
    %}
    {% func (bp *BasePage) Title() %}{%s bp.TitleStr %}{% endfunc %}
    {% func (bp *BasePage) Body(s string, n int) %}
        <b>s={%q s %}, n={%d n %}</b>
    {% endfunc %}
    {% func (bp *BasePage) Footer() %}{%s bp.FooterStr %}{% endfunc %}

    Main page implementation
    {% code
    type MainPage struct {
        // inherit from BasePage
        BasePage

        // real body for main page
        BodyStr string
    }
    %}

    Override only Body
    Title and Footer are used from BasePage.
    {% func (mp *MainPage) Body(s string, n int) %}
        <div>
            main body: {%s mp.BodyStr %}
        </div>
        <div>
            base body: {%= mp.BasePage.Body(s, n) %}
        </div>
    {% endfunc %}
    ```

    See [basicserver example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver)
    for more details.

# Performance optimization tips

  * Prefer calling `WriteFoo` instead of `Foo` when generating template output
    for `{% func Foo() %}`. This avoids unnesessary memory allocation and a copy
    for a `string` returned from `Foo()`.

  * Prefer `{%= Foo() %}` instead of `{%s= Foo() %}` when embedding
    a function template `{% func Foo() %}`. Though both approaches generate
    identical output, the first approach is optimized for speed.

  * Prefer using existing output tags instead of passing `fmt.Sprintf`
    to `{%s %}`. For instance, use `{%d num %}` instead
    of `{%s fmt.Sprintf("%d", num) %}`, because the first approach is optimized
    for speed.

  * Prefer using specific output tags instead of generic output tag
    `{%v %}`. For example, use `{%s str %}` instead of `{%v str %}`, since
    specific output tags are optimized for speed.

  * Prefer creating custom function templates instead of composing complex
    strings by hands before passing them to `{%s %}`.
    For instance, the first approach is slower than the second one:

    ```qtpl
    {% func Foo(n int) %}
        {% code
        // construct complex string
        complexStr := ""
        for i := 0; i < n; i++ {
            complexStr += fmt.Sprintf("num %d,", i)
        }
        %}
        complex string = {%s= complexStr %}
    {% endfunc %}
    ```

    ```qtpl
    {% func Foo(n int) %}
        complex string = {%= complexStr(n) %}
    {% endfunc %}

    // Wrap complexStr func into stripspace for stripping unnesessary space
    // between tags and lines.
    {% stripspace %}
    {% func complexStr(n int) %}
        {% for i := 0; i < n; i++ %}
            num{% space %}{%d i %}{% newline %}
        {% endfor %}
    {% endfunc %}
    {% endstripspace %}
    ```

  * Make sure that the `io.Writer` passed to `Write*` functions
    is [buffered](https://golang.org/pkg/bufio/#Writer).
    This will minimize the number of `write`
    [syscalls](https://en.wikipedia.org/wiki/System_call),
    which may be quite expensive.

    Note: There is no need to wrap [fasthttp.RequestCtx](https://godoc.org/github.com/valyala/fasthttp#RequestCtx)
    into [bufio.Writer](https://golang.org/pkg/bufio/#Writer), since it is already buffered.

  * [Profile](http://blog.golang.org/profiling-go-programs) your programs
    for memory allocations and fix the most demanding functions based on
    the output of `go tool pprof --alloc_objects`.

# Use cases

While the main quicktemplate purpose is generating HTML, it may be used
for generating other data too. For example, JSON and XML marshalling may
be easily implemented with quicktemplate:

```qtpl
{% code
type MarshalRow struct {
	Msg string
	N int
}

type MarshalData struct {
	Foo int
	Bar string
	Rows []MarshalRow
}
%}

// JSON marshaling
{% stripspace %}
{% func (d *MarshalData) JSON() %}
{
	"Foo": {%d d.Foo %},
	"Bar": {%q= d.Bar %},
	"Rows":[
		{% for i, r := range d.Rows %}
			{
				"Msg": {%q= r.Msg %},
				"N": {%d r.N %}
			}
			{% if i + 1 < len(d.Rows) %},{% endif %}
		{% endfor %}
	]
}
{% endfunc %}
{% endstripspace %}

// XML marshalling
{% stripspace %}
{% func (d *MarshalData) XML() %}
<MarshalData>
	<Foo>{%d d.Foo %}</Foo>
	<Bar>{%s d.Bar %}</Bar>
	<Rows>
	{% for _, r := range d.Rows %}
		<Row>
			<Msg>{%s r.Msg %}</Msg>
			<N>{%d r.N %}</N>
		</Row>
	{% endfor %}
	</Rows>
</MarshalData>
{% endfunc %}
{% endstripspace %}
```

Usually, marshalling built with quicktemplate works faster than the marshalling
implemented via standard [encoding/json](https://golang.org/pkg/encoding/json/)
and [encoding/xml](https://golang.org/pkg/encoding/xml/).
See the corresponding benchmark results:

```
go test -bench=Marshal -benchmem github.com/valyala/quicktemplate/tests
BenchmarkMarshalJSONStd1-4                	 3000000	       480 ns/op	       8 B/op	       1 allocs/op
BenchmarkMarshalJSONStd10-4               	 1000000	      1842 ns/op	       8 B/op	       1 allocs/op
BenchmarkMarshalJSONStd100-4              	  100000	     15820 ns/op	       8 B/op	       1 allocs/op
BenchmarkMarshalJSONStd1000-4             	   10000	    159327 ns/op	      59 B/op	       1 allocs/op
BenchmarkMarshalJSONQuickTemplate1-4      	10000000	       162 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONQuickTemplate10-4     	 2000000	       748 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONQuickTemplate100-4    	  200000	      6572 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalJSONQuickTemplate1000-4   	   20000	     66784 ns/op	      29 B/op	       0 allocs/op
BenchmarkMarshalXMLStd1-4                 	 1000000	      1652 ns/op	       2 B/op	       2 allocs/op
BenchmarkMarshalXMLStd10-4                	  200000	      7533 ns/op	      11 B/op	      11 allocs/op
BenchmarkMarshalXMLStd100-4               	   20000	     65763 ns/op	     195 B/op	     101 allocs/op
BenchmarkMarshalXMLStd1000-4              	    2000	    663373 ns/op	    3522 B/op	    1002 allocs/op
BenchmarkMarshalXMLQuickTemplate1-4       	10000000	       145 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLQuickTemplate10-4      	 3000000	       597 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLQuickTemplate100-4     	  300000	      5833 ns/op	       0 B/op	       0 allocs/op
BenchmarkMarshalXMLQuickTemplate1000-4    	   30000	     53000 ns/op	      32 B/op	       0 allocs/op
```

# FAQ

  * *Why is the quicktemplate syntax incompatible with [html/template](https://golang.org/pkg/html/template/)?*

    Because `html/template` syntax isn't expressive enough for `quicktemplate`.

  * *What's the difference between quicktemplate and [ego](https://github.com/benbjohnson/ego)?*

    `Ego` is similar to `quicktemplate` in the sense it converts templates into Go code.
    But it misses the following stuff, which makes `quicktemplate` so powerful
    and easy to use:

      * Defining multiple function templates in a single template file.
      * Embedding function templates inside other function templates.
      * Template interfaces, inheritance and overriding.
        See [this example](https://github.com/valyala/quicktemplate/tree/master/examples/basicserver)
        for details.
      * Top-level comments outside function templates.
      * Template packages.
      * Combining arbitrary Go files with template files in template packages.
      * Performance optimizations.

  * *What's the difference between quicktemplate and [gorazor](https://github.com/sipin/gorazor)?*

    `Gorazor` is similar to `quicktemplate` in the sense it converts templates into Go code.
    But it misses the following useful features:

      * Clear syntax insead of hard-to-understand magic stuff related
        to template arguments, template inheritance and embedding function
        templates into other templates.
      * Performance optimizations.

* *Is there a syntax highlighting for qtpl files?*

  Yes - see [this issue](https://github.com/valyala/quicktemplate/issues/19) for details.
  If you are using JetBrains products (syntax highlighting and autocomplete):
    * cd [JetBrains settings directory](https://intellij-support.jetbrains.com/hc/en-us/articles/206544519-Directories-used-by-the-IDE-to-store-settings-caches-plugins-and-logs)
    * mkdir -p filetypes && cd filetypes
    * curl https://raw.githubusercontent.com/valyala/quicktemplate/master/QuickTemplate.xml >> QuickTemplate.xml
    * Restart your IDE

* *I didn't find an answer for my question here.*

  Try exploring [these questions](https://github.com/valyala/quicktemplate/issues?q=label%3Aquestion).
