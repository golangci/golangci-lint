# godot

[![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/tetafro/godot/master/LICENSE)
[![Github CI](https://img.shields.io/github/workflow/status/tetafro/godot/Test)](https://github.com/tetafro/godot/actions?query=workflow%3ATest)
[![Go Report](https://goreportcard.com/badge/github.com/tetafro/godot)](https://goreportcard.com/report/github.com/tetafro/godot)
[![Codecov](https://codecov.io/gh/tetafro/godot/branch/master/graph/badge.svg)](https://codecov.io/gh/tetafro/godot)

Linter that checks if all top-level comments contain a period at the
end of the last sentence if needed.

[CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences) quote:

> Comments should begin with the name of the thing being described
> and end in a period

## Install and run

```sh
go get -u github.com/tetafro/godot/cmd/godot
godot ./myproject
```

## Examples

Code

```go
package math

// Sum sums two integers
func Sum(a, b int) int {
    return a + b // result
}
```

Output

```sh
Top level comment should end in a period: math/math.go:3:1
```
