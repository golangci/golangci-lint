

## GAS - Go AST Scanner

Inspects source code for security problems by scanning the Go AST.

### License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License [here](http://www.apache.org/licenses/LICENSE-2.0).

### Project status

[![Build Status](https://travis-ci.org/GoASTScanner/gas.svg?branch=master)](https://travis-ci.org/GoASTScanner/gas)
[![GoDoc](https://godoc.org/github.com/GoASTScanner/gas?status.svg)](https://godoc.org/github.com/GoASTScanner/gas)

Gas is still in alpha and accepting feedback from early adopters. We do
not consider it production ready at this time.

### Install

`$ go get github.com/GoASTScanner/gas/cmd/gas/...`

### Usage

Gas can be configured to only run a subset of rules, to exclude certain file
paths, and produce reports in different formats. By default all rules will be
run against the supplied input files. To recursively scan from the current
directory you can supply './...' as the input argument.

#### Selecting rules

By default Gas will run all rules against the supplied file paths. It is however possible to select a subset of rules to run via the '-include=' flag,
or to specify a set of rules to explicitly exclude using the '-exclude=' flag.

##### Available rules

  - G101: Look for hardcoded credentials
  - G102: Bind to all interfaces
  - G103: Audit the use of unsafe block
  - G104: Audit errors not checked
  - G105: Audit the use of math/big.Int.Exp
  - G106: Audit the use of ssh.InsecureIgnoreHostKey
  - G201: SQL query construction using format string
  - G202: SQL query construction using string concatenation
  - G203: Use of unescaped data in HTML templates
  - G204: Audit use of command execution
  - G301: Poor file permissions used when creating a directory
  - G302: Poor file permisions used with chmod
  - G303: Creating tempfile using a predictable path
  - G304: File path provided as taint input
  - G401: Detect the usage of DES, RC4, or MD5
  - G402: Look for bad TLS connection settings
  - G403: Ensure minimum RSA key length of 2048 bits
  - G404: Insecure random number source (rand)
  - G501: Import blacklist: crypto/md5
  - G502: Import blacklist: crypto/des
  - G503: Import blacklist: crypto/rc4
  - G504: Import blacklist: net/http/cgi


```
# Run a specific set of rules
$ gas -include=G101,G203,G401 ./...

# Run everything except for rule G303
$ gas -exclude=G303 ./...
```

#### Excluding files:

Gas will ignore dependencies in your vendor directory any files
that are not considered build artifacts by the compiler (so test files).

#### Annotating code

As with all automated detection tools there will be cases of false positives. In cases where Gas reports a failure that has been manually verified as being safe it is possible to annotate the code with a '#nosec' comment.

The annotation causes Gas to stop processing any further nodes within the
AST so can apply to a whole block or more granularly to a single expression.

```go

import "md5" // #nosec


func main(){

    /* #nosec */
    if x > y {
        h := md5.New() // this will also be ignored
    }

}

```

In some cases you may also want to revisit places where #nosec annotations
have been used. To run the scanner and ignore any #nosec annotations you
can do the following:

```
$ gas -nosec=true ./...
```
#### Build tags

Gas is able to pass your [Go build tags](https://golang.org/pkg/go/build/) to the analyzer.
They can be provided as a comma separated list as follows:

```
$ gas -tag debug,ignore ./...
```

### Output formats

Gas currently supports text, json, yaml, csv and JUnit XML output formats. By default
results will be reported to stdout, but can also be written to an output
file. The output format is controlled by the '-fmt' flag, and the output file is controlled by the '-out' flag as follows:

```
# Write output in json format to results.json
$ gas -fmt=json -out=results.json *.go
```
### Development

#### Prerequisites

Install dep according to the instructions here: https://github.com/golang/dep
Install the latest version of golint: https://github.com/golang/lint

#### Build

```
make
```

#### Tests

```
make test
```

#### Release Build

Gas can be released as follows:

```bash
make release VERSION=2.0.0
```

The released version of the tool is available in the `build` folder. The build information should be displayed in the usage text.

```
./build/gas-2.0.0-linux-amd64 -h

GAS - Go AST Scanner

Gas analyzes Go source code to look for common programming mistakes that
can lead to security problems.

VERSION: 2.0.0
GIT TAG: 96489ff
BUILD DATE: 2018-02-21

```

#### Docker image

You can execute a release and build the docker image as follows:

```
make image VERSION=2.0.0
```

Now you can run the gas tool in a container against your local workspace:

```
docker run -it -v <YOUR LOCAL WORKSPACE>:/workspace gas /workspace
```

#### Generate TLS rule

The configuration of TLS rule can be generated from [Mozilla's TLS ciphers recommendation](https://statics.tls.security.mozilla.org/server-side-tls-conf.json).


First you need to install the generator tool:

```
go get github.com/GoASTScanner/gas/cmd/tlsconfig/...
```

You can invoke now the `go generate` in the root of the project:

```
go generate ./...
```

This will generate the `rules/tls_config.go` file with will contain the current ciphers recommendation from Mozilla.
