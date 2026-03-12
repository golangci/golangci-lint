# Rule Documentation

## Rules accepting parameters

As [README.md](https://github.com/securego/gosec/blob/master/README.md) mentions, some rules can be configured by adding parameters to the gosec JSON config. Per rule configs are encoded as top level objects in the gosec config, with the rule ID (`Gxxx`) as the key.

Currently, the following rules accept parameters. This list is manually maintained; if you notice an omission please add it!

### G101

The hard-coded credentials rule `G101` can be configured with additional patterns, and the entropy threshold can be adjusted:

```JSON
{
    "G101": {
        "pattern": "(?i)passwd|pass|password|pwd|secret|private_key|token",
         "ignore_entropy": false,
         "entropy_threshold": "80.0",
         "per_char_threshold": "3.0",
         "truncate": "32"
    }
}
```

### G104

The unchecked error value rule `G104` can be configured with additional functions that should be permitted to be called without checking errors.

```JSON
{
    "G104": {
        "ioutil": ["WriteFile"]
    }
}
```

### G111

The HTTP Directory serving rule `G111` can be configured with a different regex for detecting potentially overly permissive servers. Note that this *replaces* the default pattern of `http\.Dir\("\/"\)|http\.Dir\('\/'\)`.

```JSON
{
    "G111": {
        "pattern": "http\\.Dir\\(\"\\\/\"\\)|http\\.Dir\\('\\\/'\\)"
    }
}

```

### G301, G302, G306, G307

The various file and directory permission checking rules can be configured with a different maximum allowable file permission.

```JSON
{
    "G301":"0o600",
    "G302":"0o600",
    "G306":"0o750",
    "G307":"0o750"
}
```
