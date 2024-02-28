## The Manual Way

- add a blank import of module inside `cmd/golangci-lint/plugins.go`
- run `go mod tidy`. (the module will be imported)
- run `make build`
- define the plugin inside the configuration `linters-settings.custom` section with the type `module`.
- run you custom version of golangci-lint

## The Automatic Way

- define your building configuration into `.mygcl.yml`
- run the command `golangci-lint custom` (`go run ./cmd/golangci-lint custom -v`)
- define the plugin inside the  `linters-settings.custom` section with the type `module`.
- run your custom version of golangci-lint

### Build Configuration Example

`.mygcl.yml`:

```yaml
version: feat/new-plugin-module
name: custom-golangci-lint
plugins:
  - module: 'github.com/bombsimon/wsl/v4'
    version: v4.2.1

  - module: 'github.com/Antonboom/testifylint'
    import: 'github.com/Antonboom/testifylint/analyzer'
    version: v1.1.2

  - module: 'golang.org/x/tools'
    import: 'golang.org/x/tools/go/analysis'
    path: /home/ldez/sources/golangci-lint/x-tools/
```

- for now, only my branch `feat/new-plugin-module` can be used as version.
- the "plugins" are not real plugins is just for the import demo
- the field `path` should point to your local copy of `golang.org/x/tools`
    - (`git clone git@github.com:golang/tools.git`)
- the field `name` is optional and can be removed, the default name is `gcl-custom`
- you can use the env var `MYGCL_KEEP_TEMP_FILES` to preserve the temp directory.
    - ex: `MYGCL_KEEP_TEMP_FILES=1 mygcl`

## Golangci-lint Configuration Example

My branch contains a fake plugin, if you run `make build`, you will be able to use it.

IMPORTANT: if you create a custom binary with the command `custom`, this fake plugin will be removed.

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
