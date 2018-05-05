rollbar
-------

`rollbar` is a Go Rollbar client that makes it easy to report errors to Rollbar
with stacktraces. Errors are sent to Rollbar asynchronously in a background
goroutine.

Because Go's `error` type doesn't include stack information from when it was set
or allocated, `rollbar` uses the stack information from where the error was
reported.

You may also want to look at:

* [stvp/roll](https://github.com/stvp/roll) - Simpler, synchronous (no
  background goroutine) with a nicer API.

Installation
=============

Standard installation to your GOPATH via go get:

```
go get github.com/stvp/rollbar
```

Documentation
=============

[API docs on godoc.org](http://godoc.org/github.com/stvp/rollbar)

Usage
=====

```go
package main

import (
  "github.com/stvp/rollbar"
)

func main() {
  rollbar.Token = "MY_TOKEN"
  rollbar.Environment = "production" // defaults to "development"

  result, err := DoSomething()
  if err != nil {
    // Error reporting
    rollbar.Error(rollbar.ERR, err)
  }

  // Message reporting
  rollbar.Message("info", "Message body goes here")

  // Block until all queued messages are sent to Rollbar.
  // You can do this in a defer() if needed.
  rollbar.Wait()
}
```

Running Tests
=============

Set up a dummy project in Rollbar and pass the access token as an environment
variable to `go test`:

    TOKEN=f0df01587b8f76b2c217af34c479f9ea go test

And verify the reported errors manually in the Rollbar dashboard.

Other Resources
===============

For best practices and more information on how to handle errors in Go, these are
some great places to get started:

* [Error Handling in Go](https://blog.golang.org/error-handling-and-go)
* [Why does Go not have exceptions?](https://golang.org/doc/faq#exceptions)
* [Defer, Panic and Recover](https://blog.golang.org/defer-panic-and-recover)
* [pkg/errors](https://github.com/pkg/errors)

Contributors
============

Thanks, all!

* @kjk
* @nazwa
* @ossareh
* @paulmach
* @Soulou
* @tike
* @tysonmote
* @marcelgruber
* @karlpatr
* @sumeet
* @dfuentes77
* @seriousben

