# MyPlugin

This is a plugin example.

## Details

The plugin is registered like that:

```go

func init() {
	register.Plugin("foo", New)
}

```

The plugin should be declared as an import inside `cmd/golangci-lint/plugins.go`

```go
package main

import _ "github.com/golangci/golangci-lint/pkg/experimental/modules/myplugin"
```

The plugin should be defined inside `.golangci.yml` like with the Go plugin system:

```yaml
linters-settings:
  custom:
    foo:
      type: "module"
      description: This is an example usage of a plugin linter.
      settings:
        message: hello

linters:
  disable-all: true
  enable:
    - foo
```
