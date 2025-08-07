---
title: Module Plugin System
weight: 1
---

> [!TIP]
> An example linter can be found at [here](https://github.com/golangci/example-plugin-module-linter).

## The Automatic Way

- Define your building configuration into `.custom-gcl.yml`.
- Run the command `golangci-lint custom` (or `golangci-lint custom -v` to have logs).
- Define the plugin inside the `linters.settings.custom` section with the type `module`.
- Run the resulting custom binary of golangci-lint (`./custom-gcl` by default).

Requirements:
- Go
- git

### Configuration Example

```yaml {filename=".custom-gcl.yml"}
version: {{< latest-version >}}
plugins:
  # a plugin from a Go proxy
  - module: 'github.com/golangci/plugin1'
    import: 'github.com/golangci/plugin1/foo'
    version: v1.0.0

  # a plugin from local source
  - module: 'github.com/golangci/plugin2'
    path: /my/local/path/plugin2
```

```yaml {filename=".golangci.yml"}
version: "2"

linters:
  default: none
  enable:
    - foo
  settings:
    custom:
      foo:
        type: "module"
        description: This is an example usage of a plugin linter.
        settings:
          message: hello
```

## The Manual Way

- Add a blank-import of your module inside `cmd/golangci-lint/plugins.go`.
- Run `go mod tidy` (the module containing the plugin will be imported).
- Run `make build`.
- Define the plugin inside the `linters.settings.custom` section with the type `module`.
- Run your custom version of golangci-lint.

### Configuration Example

```yaml {filename=".golangci.yml"}
version: "2"

linters:
  default: none
  enable:
    - foo
  settings:
    custom:
      foo:
        type: "module"
        description: This is an example usage of a plugin linter.
        settings:
          message: hello
```

## Reference

The configuration file can be validated with the JSON Schema: [custom-gcl.jsonschema.json](https://golangci-lint.run/jsonschema/custom-gcl.jsonschema.json)

```yml {filename=".custom-gcl.yml"}
{ .CustomGCLReference }
```
