Gocyclo calculates cyclomatic complexities of functions in Go source code.

The cyclomatic complexity of a function is calculated according to the
following rules:

     1 is the base complexity of a function
    +1 for each 'if', 'for', 'case', '&&' or '||'

To install, run

    $ go get github.com/fzipp/gocyclo

and put the resulting binary in one of your PATH directories if
`$GOPATH/bin` isn't already in your PATH.

Usage:

    $ gocyclo [<flag> ...] <Go file or directory> ...

Examples:

    $ gocyclo .
    $ gocyclo main.go
    $ gocyclo -top 10 src/
    $ gocyclo -over 25 docker
    $ gocyclo -avg .

The output fields for each line are:

    <complexity> <package> <function> <file:row:column>

