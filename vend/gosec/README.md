
# gosec - Go Security Checker

Inspects source code for security problems by scanning the Go AST and SSA code representation.

<img src="https://securego.io/img/gosec.png" width="320">

## Features

- **Pattern-based rules** for detecting common security issues in Go code
- **SSA-based analyzers** for type conversions, slice bounds, and crypto issues
- **Taint analysis** for tracking data flow from user input to dangerous functions (SQL injection, command injection, path traversal, SSRF, XSS, log injection)

## License

Licensed under the Apache License, Version 2.0 (the "License").
You may not use this file except in compliance with the License.
You may obtain a copy of the License [here](http://www.apache.org/licenses/LICENSE-2.0).

## Project status

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3218/badge)](https://bestpractices.coreinfrastructure.org/projects/3218)
[![Build Status](https://github.com/securego/gosec/workflows/CI/badge.svg)](https://github.com/securego/gosec/actions?query=workflows%3ACI)
[![Coverage Status](https://codecov.io/gh/securego/gosec/branch/master/graph/badge.svg)](https://codecov.io/gh/securego/gosec)
[![GoReport](https://goreportcard.com/badge/github.com/securego/gosec)](https://goreportcard.com/report/github.com/securego/gosec)
[![GoDoc](https://pkg.go.dev/badge/github.com/securego/gosec/v2)](https://pkg.go.dev/github.com/securego/gosec/v2)
[![Docs](https://readthedocs.org/projects/docs/badge/?version=latest)](https://securego.io/)
[![Downloads](https://img.shields.io/github/downloads/securego/gosec/total.svg)](https://github.com/securego/gosec/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/securego/gosec.svg)](https://hub.docker.com/r/securego/gosec/tags)
[![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](http://securego.slack.com)
[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)

## Install

### CI Installation

```bash
# binary will be $(go env GOPATH)/bin/gosec
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin vX.Y.Z

# or install it into ./bin/
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s vX.Y.Z

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s vX.Y.Z

# If you want to use the checksums provided on the "Releases" page
# then you will have to download a tar.gz file for your operating system instead of a binary file
wget https://github.com/securego/gosec/releases/download/vX.Y.Z/gosec_vX.Y.Z_OS.tar.gz

# The file will be in the current folder where you run the command
# and you can check the checksum like this
echo "<check sum from the check sum file>  gosec_vX.Y.Z_OS.tar.gz" | sha256sum -c -

gosec --help
```

### GitHub Action

You can run `gosec` as a GitHub action as follows:

```yaml
name: Run Gosec
on:
  push:
    branches:
      - master
  pull_request:
    branches:
      - master
jobs:
  tests:
    runs-on: ubuntu-latest
    env:
      GO111MODULE: on
    steps:
      - name: Checkout Source
        uses: actions/checkout@v3
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: ./...
```

#### Scanning Projects with Private Modules

If your project imports private Go modules, you need to configure authentication so that `gosec` can fetch the dependencies. Set the following environment variables in your workflow:

- `GOPRIVATE`: A comma-separated list of module path prefixes that should be considered private (e.g., `github.com/your-org/*`).
- `GITHUB_AUTHENTICATION_TOKEN`: A GitHub token with read access to your private repositories.

```yaml
name: Run Gosec
on:
  push:
    branches:
      - master
  pull_request:
    branches:
      - master
jobs:
  tests:
    runs-on: ubuntu-latest
    env:
      GO111MODULE: on
      GOPRIVATE: github.com/your-org/*
      GITHUB_AUTHENTICATION_TOKEN: ${{ secrets.PRIVATE_REPO_TOKEN }}
    steps:
      - name: Checkout Source
        uses: actions/checkout@v3
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: ./...
```

### Integrating with code scanning

You can [integrate third-party code analysis tools](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/integrating-with-code-scanning) with GitHub code scanning by uploading data as SARIF files.

The workflow shows an example of running the `gosec` as a step in a GitHub action workflow which outputs the `results.sarif` file. The workflow then uploads the `results.sarif` file to GitHub using the `upload-sarif` action.

```yaml
name: "Security Scan"

# Run workflow each time code is pushed to your repository and on a schedule.
# The scheduled workflow runs every at 00:00 on Sunday UTC time.
on:
  push:
  schedule:
  - cron: '0 0 * * 0'

jobs:
  tests:
    runs-on: ubuntu-latest
    env:
      GO111MODULE: on
    steps:
      - name: Checkout Source
        uses: actions/checkout@v3
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          # we let the report trigger content trigger a failure using the GitHub Security features.
          args: '-no-fail -fmt sarif -out results.sarif ./...'
      - name: Upload SARIF file
        uses: github/codeql-action/upload-sarif@v2
        with:
          # Path to SARIF file relative to the root of the repository
          sarif_file: results.sarif
```

### Go Analysis

The `goanalysis` package provides a [`golang.org/x/tools/go/analysis.Analyzer`](https://pkg.go.dev/golang.org/x/tools/go/analysis) for integration with tools that support the standard Go analysis interface, such as Bazel's [nogo](https://github.com/bazelbuild/rules_go/blob/master/go/nogo.rst) framework:

```starlark
nogo(
    name = "nogo",
    deps = [
        "@com_github_securego_gosec_v2//goanalysis",
        # add more analyzers as needed
    ],
    visibility = ["//visibility:public"],
)
```

### Local Installation

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

## Usage

Gosec can be configured to only run a subset of rules, to exclude certain file
paths, and produce reports in different formats. By default all rules will be
run against the supplied input files. To recursively scan from the current
directory you can supply `./...` as the input argument.

### Available rules

- G101: Look for hard coded credentials
- G102: Bind to all interfaces
- G103: Audit the use of unsafe block
- G104: Audit errors not checked
- G106: Audit the use of ssh.InsecureIgnoreHostKey
- G107: Url provided to HTTP request as taint input
- G108: Profiling endpoint automatically exposed on /debug/pprof
- G109: Potential Integer overflow made by strconv.Atoi result conversion to int16/32
- G110: Potential DoS vulnerability via decompression bomb
- G111: Potential directory traversal
- G112: Potential slowloris attack
- G114: Use of net/http serve function that has no support for setting timeouts
- G115: Potential integer overflow when converting between integer types
- G116: Detect Trojan Source attacks using bidirectional Unicode control characters
- G117: Potential exposure of secrets via JSON marshaling
- G201: SQL query construction using format string
- G202: SQL query construction using string concatenation
- G203: Use of unescaped data in HTML templates
- G204: Audit use of command execution
- G301: Poor file permissions used when creating a directory
- G302: Poor file permissions used with chmod
- G303: Creating tempfile using a predictable path
- G304: File path provided as taint input
- G305: File traversal when extracting zip/tar archive
- G306: Poor file permissions used when writing to a new file
- G307: Poor file permissions used when creating a file with os.Create
- G401: Detect the usage of MD5 or SHA1
- G402: Look for bad TLS connection settings
- G403: Ensure minimum RSA key length of 2048 bits
- G404: Insecure random number source (rand)
- G405: Detect the usage of DES or RC4
- G406: Detect the usage of MD4 or RIPEMD160
- G407: Detect the usage of hardcoded Initialization Vector(IV)/Nonce
- G501: Import blocklist: crypto/md5
- G502: Import blocklist: crypto/des
- G503: Import blocklist: crypto/rc4
- G504: Import blocklist: net/http/cgi
- G505: Import blocklist: crypto/sha1
- G506: Import blocklist: golang.org/x/crypto/md4
- G507: Import blocklist: golang.org/x/crypto/ripemd160
- G601: Implicit memory aliasing of items from a range statement (only for Go 1.21 or lower)
- G602: Slice access out of bounds
- G701: SQL injection via taint analysis
- G702: Command injection via taint analysis
- G703: Path traversal via taint analysis
- G704: SSRF via taint analysis
- G705: XSS via taint analysis
- G706: Log injection via taint analysis

### Retired rules

- G105: Audit the use of math/big.Int.Exp - [CVE is fixed](https://github.com/golang/go/issues/15184)
- G113: Usage of Rat.SetString in math/big with an overflow (CVE-2022-23772). This affected Go <1.16.14 and Go <1.17.7, which are no longer supported by gosec. 
- G307: Deferring a method which returns an error - causing more inconvenience than fixing a security issue, despite the details from this [blog post](https://www.joeshaw.org/dont-defer-close-on-writable-files/)

### Selecting rules

By default, gosec will run all rules against the supplied file paths. It is however possible to select a subset of rules to run via the `-include=` flag,
or to specify a set of rules to explicitly exclude using the `-exclude=` flag.

```bash
# Run a specific set of rules
$ gosec -include=G101,G203,G401 ./...

# Run everything except for rule G303
$ gosec -exclude=G303 ./...
```

### CWE Mapping

Every issue detected by `gosec` is mapped to a [CWE (Common Weakness Enumeration)](http://cwe.mitre.org/data/index.html) which describes in more generic terms the vulnerability. The exact mapping can be found  [here](https://github.com/securego/gosec/blob/master/issue/issue.go#L50).

### Configuration

A number of global settings can be provided in a configuration file as follows:

```JSON
{
    "global": {
        "nosec": "enabled",
        "audit": "enabled"
    }
}
```

- `nosec`: this setting will overwrite all `#nosec` directives defined throughout the code base
- `audit`: runs in audit mode which enables addition checks that for normal code analysis might be too nosy

```bash
# Run with a global configuration file
$ gosec -conf config.json .
```

### Path-Based Rule Exclusions

Large repositories with multiple components may need different security rules
for different paths. Use `exclude-rules` to suppress specific rules for specific
paths.

**Configuration File:**
```json
{
  "exclude-rules": [
    {
      "path": "cmd/.*",
      "rules": ["G204", "G304"]
    },
    {
      "path": "scripts/.*",
      "rules": ["*"]
    }
  ]
}
```

**CLI Flag:**
```bash
# Exclude G204 and G304 from cmd/ directory
gosec --exclude-rules="cmd/.*:G204,G304" ./...

# Exclude all rules from scripts/ directory  
gosec --exclude-rules="scripts/.*:*" ./...

# Multiple exclusions
gosec --exclude-rules="cmd/.*:G204,G304;test/.*:G101" ./...
```

| Field | Type | Description |
|-------|------|-------------|
| `path` | string (regex) | Regular expression matched against file paths |
| `rules` | []string | Rule IDs to exclude. Use `*` to exclude all rules |

#### Rule Configuration

Some rules accept configuration flags as well; these flags are documented in [RULES.md](https://github.com/securego/gosec/blob/master/RULES.md).

#### Go version

Some rules require a specific Go version which is retrieved from the Go module file present in the project. If this version cannot be found, it will fallback to Go runtime version.

The Go module version is parsed using the `go list` command which in some cases might lead to performance degradation. In this situation, the go module version can be easily provided by setting the environment variable `GOSECGOVERSION=go1.21.1`.

### Dependencies

gosec will fetch automatically the dependencies of the code which is being analyzed when go module is turned on (e.g.`GO111MODULE=on`). If this is not the case,
the dependencies need to be explicitly downloaded by running the `go get -d` command before the scan.

### Excluding test files and folders

gosec will ignore test files across all packages and any dependencies in your vendor directory.

The scanning of test files can be enabled with the following flag:

```bash
gosec -tests ./...
```

Also additional folders can be excluded as follows:

```bash
 gosec -exclude-dir=rules -exclude-dir=cmd ./...
```

### Excluding generated files

gosec can ignore generated go files with default generated code comment.

```
// Code generated by some generator DO NOT EDIT.
```

```bash
gosec -exclude-generated ./...
```

### Auto fixing vulnerabilities

gosec can suggest fixes based on AI recommendation. It will call an AI API to receive a suggestion for a security finding.

You can enable this feature by providing the following command line arguments:

- `ai-api-provider`: the name of the AI API provider. Supported providers:
  - **Gemini**: `gemini-2.5-pro`, `gemini-2.5-flash`, `gemini-2.5-flash-lite`, `gemini-2.0-flash`, `gemini-2.0-flash-lite` (default)
  - **Claude**: `claude-sonnet-4-0` (default), `claude-opus-4-0`, `claude-opus-4-1`, `claude-sonnet-3-7`
  - **OpenAI**: `gpt-4o` (default), `gpt-4o-mini`
  - **Custom OpenAI-compatible**: Any custom model name (requires `ai-base-url`)
- `ai-api-key` or set the environment variable `GOSEC_AI_API_KEY`: the key to access the AI API
  - For Gemini, you can create an API key following [these instructions](https://ai.google.dev/gemini-api/docs/api-key)
  - For Claude, get your API key from [Anthropic Console](https://console.anthropic.com/)
  - For OpenAI, get your API key from [OpenAI Platform](https://platform.openai.com/api-keys)
- `ai-base-url`: (optional) custom base URL for OpenAI-compatible APIs (e.g., Azure OpenAI, LocalAI, Ollama)
- `ai-skip-ssl`: (optional) skip SSL certificate verification for AI API (useful for self-signed certificates)

**Examples:**

```bash
# Using Gemini
gosec -ai-api-provider="gemini-2.0-flash" -ai-api-key="your_key" ./...

# Using Claude
gosec -ai-api-provider="claude-sonnet-4-0" -ai-api-key="your_key" ./...

# Using OpenAI
gosec -ai-api-provider="gpt-4o" -ai-api-key="your_key" ./...

# Using Azure OpenAI
gosec -ai-api-provider="gpt-4o" \
  -ai-api-key="your_azure_key" \
  -ai-base-url="https://your-resource.openai.azure.com/openai/deployments/your-deployment" \
  ./...

# Using local Ollama with custom model
gosec -ai-api-provider="llama3.2" \
  -ai-base-url="http://localhost:11434/v1" \
  ./...

# Using self-signed certificate API
gosec -ai-api-provider="custom-model" \
  -ai-api-key="your_key" \
  -ai-base-url="https://internal-api.company.com/v1" \
  -ai-skip-ssl \
  ./...
```

### Annotating code

As with all automated detection tools, there will be cases of false positives.
In cases where gosec reports a failure that has been manually verified as being safe,
it is possible to annotate the code with a comment that starts with `#nosec`.

The `#nosec` comment should have the format `#nosec [RuleList] [-- Justification]`.

The `#nosec` comment needs to be placed on the line where the warning is reported.

```go
func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		},
	}

	client := &http.Client{Transport: tr}
	_, err := client.Get("https://go.dev/")
	if err != nil {
		fmt.Println(err)
	}
}
```

When a specific false positive has been identified and verified as safe, you may
wish to suppress only that single rule (or a specific set of rules) within a section of code,
while continuing to scan for other problems. To do this, you can list the rule(s) to be suppressed within
the `#nosec` annotation, e.g: `/* #nosec G401 */` or `//#nosec G201 G202 G203`

You could put the description or justification text for the annotation. The
justification should be after the rule(s) to suppress and start with two or
more dashes, e.g: `//#nosec G101 G102 -- This is a false positive`

Alternatively, gosec also supports the `//gosec:disable` directive, which functions similar to `#nosec`:

```go
//gosec:disable G101 -- This is a false positive
```

In some cases you may also want to revisit places where `#nosec` or `//gosec:disable` annotations
have been used. To run the scanner and ignore any `#nosec` annotations you
can do the following:

```bash
gosec -nosec=true ./...
```

### Tracking suppressions

As described above, we could suppress violations externally (using `-include`/
`-exclude`) or inline (using `#nosec` annotations) in gosec. This suppression
inflammation can be used to generate corresponding signals for auditing
purposes.

We could track suppressions by the `-track-suppressions` flag as follows:

```bash
gosec -track-suppressions -exclude=G101 -fmt=sarif -out=results.sarif ./...
```

- For external suppressions, gosec records suppression info where `kind` is
`external` and `justification` is a certain sentence "Globally suppressed".
- For inline suppressions, gosec records suppression info where `kind` is
`inSource` and `justification` is the text after two or more dashes in the
comment.

**Note:** Only SARIF and JSON formats support tracking suppressions.

### Build tags

gosec is able to pass your [Go build tags](https://pkg.go.dev/go/build/) to the analyzer.
They can be provided as a comma separated list as follows:

```bash
gosec -tags debug,ignore ./...
```

### Output formats

gosec currently supports `text`, `json`, `yaml`, `csv`, `sonarqube`, `JUnit XML`, `html` and `golint` output formats. By default
results will be reported to stdout, but can also be written to an output
file. The output format is controlled by the `-fmt` flag, and the output file is controlled by the `-out` flag as follows:

```bash
# Write output in json format to results.json
$ gosec -fmt=json -out=results.json *.go
```

Results will be reported to stdout as well as to the provided output file by `-stdout` flag. The `-verbose` flag overrides the
output format when stdout the results while saving them in the output file
```bash
# Write output in json format to results.json as well as stdout
$ gosec -fmt=json -out=results.json -stdout *.go

# Overrides the output format to 'text' when stdout the results, while writing it to results.json
$ gosec -fmt=json -out=results.json -stdout -verbose=text *.go
```

**Note:** gosec generates the [generic issue import format](https://docs.sonarqube.org/latest/analysis/generic-issue/) for SonarQube, and a report has to be imported into SonarQube using `sonar.externalIssuesReportPaths=path/to/gosec-report.json`.

## Development

[CONTRIBUTING.md](https://github.com/securego/gosec/blob/master/CONTRIBUTING.md) contains detailed information about adding new rules to gosec.

### Creating Taint Analysis Rules

gosec supports taint analysis to track data flow from untrusted sources (user input) to dangerous sinks (functions that could cause security vulnerabilities). The taint analysis rules detect issues like SQL injection, command injection, path traversal, SSRF, XSS, and log injection.

#### Creating a New Taint Rule

To create a new taint analysis rule:

1. **Create the analyzer file** in `analyzers/` (e.g., `analyzers/newvuln.go`) with both the configuration and analyzer:

```go
package analyzers

import (
    "golang.org/x/tools/go/analysis"
    "github.com/securego/gosec/v2/taint"
)

// NewVulnerability returns a configuration for detecting your vulnerability
func NewVulnerability() taint.Config {
    return taint.Config{
        Sources: []taint.Source{
            {Package: "net/http", Name: "Request", Pointer: true},
            {Package: "os", Name: "Args"},
        },
        Sinks: []taint.Sink{
            {Package: "dangerous/package", Method: "DangerousFunc"},
            {Package: "another/pkg", Receiver: "Type", Method: "Method", Pointer: true},
            {Package: "database/sql", Receiver: "DB", Method: "Query", Pointer: true, CheckArgs: []int{1}},
        },
    }
}

func newNewVulnAnalyzer(id string, description string) *analysis.Analyzer {
    config := NewVulnerability()
    rule := NewVulnerabilityRule  // Define this as a variable in the same file
    rule.ID = id
    rule.Description = description
    return taint.NewGosecAnalyzer(&rule, &config)
}
```

**Note**: Each taint analyzer keeps its configuration function in the same file as the analyzer. For examples, see:
- `analyzers/sqlinjection.go` - SQL injection detection (G701)
- `analyzers/commandinjection.go` - Command injection detection (G702)
- `analyzers/pathtraversal.go` - Path traversal detection (G703)

2. **Register the analyzer** in `analyzers/analyzerslist.go`:

```go
var defaultAnalyzers = []AnalyzerDefinition{
    // ... existing analyzers ...
    {"G7XX", "Description of vulnerability", newNewVulnAnalyzer},
}
```

3. **Add test samples** in `testutils/g7xx_samples.go`:

```go
package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG7XX - Description of vulnerability
var SampleCodeG7XX = []CodeSample{
    {[]string{`
package main

import (
    "dangerous/package"
    "net/http"
)

func handler(r *http.Request) {
    input := r.URL.Query().Get("param")
    dangerous.DangerousFunc(input)  // Should detect
}
`}, 1, gosec.NewConfig()},
    {[]string{`
package main

import (
    "dangerous/package"
)

func safeHandler() {
    // Safe - no user input
    dangerous.DangerousFunc("constant")
}
`}, 0, gosec.NewConfig()},
}
```

Then add the test case to `analyzers/analyzers_test.go`:

```go
It("should detect your new vulnerability", func() {
    runner("G7XX", testutils.SampleCodeG7XX)
})
```

#### Source and Sink Configuration

**Sources** define where tainted (untrusted) data originates:
- `Package`: The import path (e.g., `"net/http"`)
- `Name`: The type or function name (e.g., `"Request"`)
- `Pointer`: Set to `true` if it's a pointer type (e.g., `*http.Request`)

**Sinks** define dangerous functions that should not receive tainted data:
- `Package`: The import path (e.g., `"database/sql"`)
- `Receiver`: The type name for methods (e.g., `"DB"`), or empty for package functions
- `Method`: The function or method name (e.g., `"Query"`)
- `Pointer`: Set to `true` if the receiver is a pointer type
- `CheckArgs`: Optional list of argument indices to check (e.g., `[]int{1}` to check only the second argument). If omitted, all arguments are checked. Useful when some arguments are safe (like prepared statement parameters) or should be ignored (like writer arguments in `fmt.Fprintf`)

**Example with CheckArgs:**
```go
// For SQL methods, Args[0] is the receiver (*sql.DB), Args[1] is the query string
// Only check the query string; prepared statement parameters (Args[2+]) are safe
{Package: "database/sql", Receiver: "DB", Method: "Query", Pointer: true, CheckArgs: []int{1}}

// For fmt.Fprintf, Args[0] is the writer (os.Stderr), Args[1+] are format and data
// Skip the writer argument, only check format string and data arguments
{Package: "fmt", Method: "Fprintf", CheckArgs: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
```

#### Common Taint Sources

| Source Type | Package | Type/Method | Pointer |
|-------------|---------|-------------|---------|
| HTTP Request | `net/http` | `Request` | `true` |
| Query Parameters | `net/http` | `Request.URL.Query()` | - |
| Form Data | `net/http` | `Request.FormValue()` | - |
| Headers | `net/http` | `Request.Header` | - |
| Command Line Args | `os` | `Args` | `false` |
| Environment Variables | `os` | `Getenv` | `false` |
| File Content | `bufio` | `Reader` | `true` |

### Build

You can build the binary with:

```bash
make
```

### Note on Sarif Types Generation

Install the tool with :

```bash
go get -u github.com/a-h/generate/cmd/schema-generate
```

Then generate the types with :

```bash
schema-generate -i sarif-schema-2.1.0.json -o mypath/types.go
```

Most of the MarshallJSON/UnmarshalJSON are removed except the one for PropertyBag which is handy to inline the additional properties. The rest can be removed.
The URI,ID, UUID, GUID were renamed so it fits the Go convention defined [here](https://github.com/golang/lint/blob/master/lint.go#L700)

### Tests

You can run all unit tests using:

```bash
make test
```

### Release

You can create a release by tagging the version as follows:

``` bash
git tag v1.0.0 -m "Release version v1.0.0"
git push origin v1.0.0
```

The GitHub [release workflow](.github/workflows/release.yml) triggers immediately after the tag is pushed upstream. This flow will
release the binaries using the [goreleaser](https://goreleaser.com/actions/) action and then it will build and publish the docker image into Docker Hub.

The released artifacts are signed using [cosign](https://docs.sigstore.dev/). You can use the public key from [cosign.pub](cosign.pub)
file to verify the signature of docker image and binaries files.

The docker image signature can be verified with the following command:
```
cosign verify --key cosign.pub securego/gosec:<TAG>
```

The binary files signature can be verified with the following command:
```
cosign verify-blob --key cosign.pub --signature gosec_<VERSION>_darwin_amd64.tar.gz.sig  gosec_<VERSION>_darwin_amd64.tar.gz
```

### Docker image

You can also build locally the docker image by using the command:

```bash
make image
```

You can run the `gosec` tool in a container against your local Go project. You only have to mount the project
into a volume as follows:

```bash
docker run --rm -it -w /<PROJECT>/ -v <YOUR PROJECT PATH>/<PROJECT>:/<PROJECT> securego/gosec /<PROJECT>/...
```

**Note:** the current working directory needs to be set with `-w` option in order to get successfully resolved the dependencies from go module file

### Generate TLS rule

The configuration of TLS rule can be generated from [Mozilla's TLS ciphers recommendation](https://statics.tls.security.mozilla.org/server-side-tls-conf.json).

First you need to install the generator tool:

```bash
go get github.com/securego/gosec/v2/cmd/tlsconfig/...
```

You can invoke now the `go generate` in the root of the project:

```bash
go generate ./...
```

This will generate the `rules/tls_config.go` file which will contain the current ciphers recommendation from Mozilla.

## Who is using gosec?

This is a [list](USERS.md) with some of the gosec's users.

## Sponsors

Support this project by becoming a sponsor. Your logo will show up here with a link to your website

<a href="https://github.com/mercedes-benz" target="_blank"><img src="https://avatars.githubusercontent.com/u/34240465?s=80&v=4"></a>
