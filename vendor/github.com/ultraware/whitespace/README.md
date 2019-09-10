# Whitespace linter

Whitespace is a linter that checks for unnecessary newlines at the start and end of functions, if, for, etc.

Example code:

```go
package main

import "fmt"

func main() {

	fmt.Println("Hello world")
}
```

Reults in:

```
$ whitespace .
main.go:6:unnecessary newline
```

## Installation guide

```bash
go get git.ultraware.nl/NiseVoid/whitespace
```

### Gometalinter

You can add whitespace to gometalinter and enable it.

`.gometalinter.json`:

```json
{
	"Linters": {
		"whitespace": "whitespace:PATH:LINE:MESSAGE"
	},

	"Enable": [
		"whitespace"
	]
}
```

commandline:

```bash
gometalinter --linter "whitespace:whitespace:PATH:LINE:MESSAGE" --enable "whitespace"
```
