# Depguard

Go linter that checks package imports are in a list of acceptable packages. It
supports a white list and black list option and only checks the prefixes of the
import path. This allows you to allow imports from a whole organization or only
allow specific packages within a repository.

## Install

```bash
go get -u github.com/OpenPeeDeeP/depguard
```

## Config

By default, Depguard looks for a file named `.depguard.json` in the current
current working directory. If it is somewhere else, pass in the `-c` flag with
the location of your configuration file.

The following is an example configuration file.

```json
{
  "type": "whitelist",
  "packages": [
    "github.com/OpenPeeDeeP/depguard"
  ],
  "includeGoRoot": true
}
```

- `type` can be either `whitelist` or `blacklist`. This check is case insensitive.
If not specified the default is `blacklist`.
- `packages` is a list of packages for the list type specified.
- Set `includeGoRoot` to true if you want to check the list against standard lib.
If not specified the default is false.

## Gometalinter

The binary installation of this linter can be used with
[Gometalinter](github.com/alecthomas/gometalinter).

If you use a configuration file for Gometalinter then the following will need to
be added to your configuration file.

```json
{
  "linters": {
    "depguard": {
      "command": "depguard -c path/to/config.json",
      "pattern": "PATH:LINE:COL:MESSAGE",
      "installFrom": "github.com/OpenPeeDeeP/depguard",
      "isFast": true,
      "partitionStrategy": "packages"
    }
  }
}
```

If you prefer the command line way the following will work for you as well.

```bash
gometalinter --linter='depguard:depguard -c path/to/config.json:PATH:LINE:COL:MESSAGE'
```

## Golangci-lint

This linter was built with
[Golangci-lint](https://github.com/golangci/golangci-lint) in mind. I have opened
up a [Pull Request](github.com/golangci/golangci-lint/pull/47) into GolangCI-Lint
and will hopefully get that merged in soon.
