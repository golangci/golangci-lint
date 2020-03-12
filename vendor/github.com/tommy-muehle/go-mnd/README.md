# go-mnd - Magic number detector for Golang

A vet analyzer to detect magic numbers.

> **What is a magic number?**  
> A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. It's considered a bad programming practice to use numbers directly in any source code without an explanation. It makes programs harder to read, understand, and maintain.

## Project status

[![Build Status](https://travis-ci.org/tommy-muehle/go-mnd.svg?branch=master)](https://travis-ci.org/tommy-muehle/go-mnd)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommy-muehle/go-mnd)](https://goreportcard.com/report/github.com/tommy-muehle/go-mnd)

## Install

This analyzer requires Golang in version >= 1.12 because it's depends on the **go/analysis** API.

```
go get github.com/tommy-muehle/go-mnd/cmd/mnd
```

To install with [Homebrew](https://brew.sh/), run:

```
brew tap tommy-muehle/tap && brew install tommy-muehle/tap/mnd
```

On Windows download the [latest release](https://github.com/tommy-muehle/go-mnd/releases).

## Usage

[![asciicast](https://asciinema.org/a/231021.svg)](https://asciinema.org/a/231021)

```
go vet -vettool $(which mnd) ./...
```

or directly

```
mnd ./...
```

The ```-checks``` option let's you define a comma separated list of checks.

The ```-ignored-numbers``` option let's you define a comma separated list of numbers to ignore.

The ```-excludes``` option let's you define a comma separated list of regexp patterns to exclude.

## Checks

By default this detector analyses arguments, assigns, cases, conditions, operations and return statements.

* argument

```
t := http.StatusText(200)
```

* assign

```
c := &http.Client{
    Timeout: 5 * time.Second,
}
```

* case

```
switch x {
    case 3:
}
```

* condition

```
if x > 7 {
}
```

* operation

```
var x, y int
y = 10 * x
```

* return

```
return 3
```

## Excludes

By default the numbers 0 and 1 as well as test files are excluded! 

### Further known excludes

The function "Date" in the "Time" package.

```
t := time.Date(2017, time.September, 26, 12, 13, 14, 0, time.UTC)
```

Additional custom excludes can be defined via option flag.

## License

The MIT License (MIT). Please see [LICENSE](LICENSE) for more information.
