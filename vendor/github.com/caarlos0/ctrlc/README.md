# ctrlc

CTRL-C is a Go library that provides an easy way of having a task that
is context-aware and deals with SIGINT and SIGTERM signals.

## Usage

```go
package main

import "context"
import "github.com/caarlos0/ctrlc"

func main() {
    ctx, cancel := context.WithTimeout(context.Backgroud(), time.Second)
    defer cancel()
    ctrlc.Default.Run(ctx, func() error {
        // do something
        return nil
    })
}
```

