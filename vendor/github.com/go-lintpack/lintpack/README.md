## Quick start / Installation / Usage

Install `lintpack`:

```bash
go get -v -u github.com/lintpack/lintpack/...
```

Install checkers from [go-critic/checkers](https://github.com/go-critic/checkers):

```bash
# You'll need to have sources under your Go workspace first:
go get -v -u github.com/go-critic/checkers
# Now build a linter that includes all checks from that package:
lintpack build -o gocritic github.com/go-critic/checkers
# Executable gocritic is created and can be used as a standalone linter.
```

Produced binary includes basic help as well as supported checks documentation.

So, the process is simple:

* Get the `lintpack` linter builder
* Build linter from checks implemented in different repos, by various vendors
