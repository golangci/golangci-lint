formatters:
  settings:
    gci:
      # Section configuration to compare against.
      # Section names are case-insensitive and may contain parameters in ().
      # The default order of sections is `standard > default > custom > blank > dot > alias > localmodule`.
      # If `custom-order` is `true`, it follows the order of `sections` option.
      # Default: ["standard", "default"]
      sections:
        - standard # Standard section: captures all standard packages.
        - default # Default section: contains all imports that could not be matched to another section type.
        - prefix(github.com/org/project) # Custom section: groups all imports with the specified Prefix.
        - blank # Blank section: contains all blank imports. This section is not present unless explicitly enabled.
        - dot # Dot section: contains all dot imports. This section is not present unless explicitly enabled.
        - alias # Alias section: contains all alias imports. This section is not present unless explicitly enabled.
        - localmodule # Local module section: contains all local packages. This section is not present unless explicitly enabled.
      # Checks that no inline comments are present.
      # Default: false
      no-inline-comments: true
      # Checks that no prefix comments (comment lines above an import) are present.
      # Default: false
      no-prefix-comments: true
      # Enable custom order of sections.
      # If `true`, make the section order the same as the order of `sections`.
      # Default: false
      custom-order: true
      # Drops lexical ordering for custom sections.
      # Default: false
      no-lex-order: true
