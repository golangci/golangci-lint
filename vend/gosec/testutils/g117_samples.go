// testutils/g117_samples.go
package testutils

import "github.com/securego/gosec/v2"

var SampleCodeG117 = []CodeSample{
	// Positive: match on field name (default JSON key)
	{[]string{`
package main

type Config struct {
	Password string
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

type Config struct {
	APIKey *string ` + "`json:\"api_key\"`" + `
}
`}, 1, gosec.NewConfig()},

	{[]string{`
package main

type Config struct {
	PrivateKey []byte ` + "`json:\"private_key\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: match on field name (explicit non-sensitive JSON key)
	{[]string{`
package main

type Config struct {
	Password string ` + "`json:\"text_field\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: match on JSON key (non-sensitive field name)
	{[]string{`
package main

type Config struct {
	SafeField string ` + "`json:\"api_key\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: match on both
	{[]string{`
package main

type Config struct {
	Token string ` + "`json:\"auth_token\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: snake/hyphen variants in JSON key
	{[]string{`
package main

type Config struct {
	Key string ` + "`json:\"access-key\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: empty json tag part falls back to field name
	{[]string{`
package main

type Config struct {
	Secret string ` + "`json:\",omitempty\"`" + `
}
`}, 1, gosec.NewConfig()},

	// Positive: plural forms
	{[]string{`
package main

type Config struct {
	ApiTokens []string
}
`}, 1, gosec.NewConfig()},

	{[]string{`
package main

type Config struct {
	RefreshTokens []string ` + "`json:\"refresh_tokens\"`" + `
}
`}, 1, gosec.NewConfig()},

	{[]string{`
package main

type Config struct {
	AccessTokens []*string
}
`}, 1, gosec.NewConfig()},

	{[]string{`
package main

type Config struct {
	CustomSecret string ` + "`json:\"my_custom_secret\"`" + `
}
`}, 1, func() gosec.Config {
		cfg := gosec.NewConfig()
		cfg.Set("G117", map[string]interface{}{
			"pattern": "(?i)custom[_-]?secret",
		})
		return cfg
	}()},

	// Negative: json:"-" (omitted)
	{[]string{`
package main

type Config struct {
	Password string ` + "`json:\"-\"`" + `
}
`}, 0, gosec.NewConfig()},

	// Negative: both field name and JSON key non-sensitive
	{[]string{`
package main

type Config struct {
	UserID string ` + "`json:\"user_id\"`" + `
}
`}, 0, gosec.NewConfig()},

	// Negative: unexported field
	{[]string{`
package main

type Config struct {
	password string
}
`}, 0, gosec.NewConfig()},

	// Negative: non-sensitive type (int) even with "token"
	{[]string{`
package main

type Config struct {
	MaxTokens int
}
`}, 0, gosec.NewConfig()},

	// Negative: non-secret plural slice (common FP like redaction placeholders)
	{[]string{`
package main

type Config struct {
	RedactionTokens []string ` + "`json:\"redactionTokens,omitempty\"`" + `
}
`}, 0, gosec.NewConfig()},

	// Negative: grouped fields, only one sensitive (should still flag the sensitive one)
	// Note: we expect 1 issue (for the sensitive field)
	{[]string{`
package main

type Config struct {
	Safe, Password string
}
`}, 1, gosec.NewConfig()},

	// Suppression: trailing line comment
	{[]string{`
package main

type Config struct {
	Password string // #nosec G117
}
`}, 0, gosec.NewConfig()},

	// Suppression: line comment above field
	{[]string{`
package main

type Config struct {
	// #nosec G117 -- false positive
	Password string
}
`}, 0, gosec.NewConfig()},

	// Suppression: trailing with justification
	{[]string{`
package main

type Config struct {
	APIKey string ` + "`json:\"api_key\"`" + ` // #nosec G117 -- public key
}
`}, 0, gosec.NewConfig()},
}
