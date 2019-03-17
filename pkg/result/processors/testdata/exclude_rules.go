package testdata

//go:generate --long line --with a --lot of --arguments --that we --would like --to exclude --from lll --issues --by exclude-rules

// long line that we don't want to exclude from lll issues. Use the similar pattern: go:generate. This line should be reported by lll
