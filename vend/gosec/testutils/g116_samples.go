package testutils

import "github.com/securego/gosec/v2"

// #nosec - This file intentionally contains bidirectional Unicode characters
// for testing trojan source detection.　The G116 rule scans the entire file content　(not just AST nodes)
// because trojan source attacks work by manipulating　visual representation of code through bidirectional
// text control characters, which can appear in comments, strings or anywhere in the source file.
// Without this #nosec exclusion, gosec would detect these test samples as actual vulnerabilities.
var (
	// SampleCodeG116 - TrojanSource code snippets
	SampleCodeG116 = []CodeSample{
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// This comment contains bidirectional unicode: access\u202e\u2066 granted\u2069\u202d\n\tisAdmin := false\n\tfmt.Println(\"Access status:\", isAdmin)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Trojan source with RLO character\n\taccessLevel := \"user\"\n\t// Actually assigns \"nimda\" due to bidi chars: accessLevel = \"\u202enimda\"\n\tif accessLevel == \"admin\" {\n\t\tfmt.Println(\"Access granted\")\n\t}\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// String with bidirectional override\n\tusername := \"admin\u202e \u2066Check if admin\u2069 \u2066\"\n\tpassword := \"secret\"\n\tfmt.Println(username, password)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains LRI (Left-to-Right Isolate) U+2066\n\tcomment := \"Safe comment \u2066with hidden text\u2069\"\n\tfmt.Println(comment)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains RLI (Right-to-Left Isolate) U+2067\n\tmessage := \"Normal text \u2067hidden\u2069\"\n\tfmt.Println(message)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains FSI (First Strong Isolate) U+2068\n\ttext := \"Text with \u2068hidden content\u2069\"\n\tfmt.Println(text)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains LRE (Left-to-Right Embedding) U+202A\n\tembedded := \"Text with \u202aembedded\u202c content\"\n\tfmt.Println(embedded)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains RLE (Right-to-Left Embedding) U+202B\n\trtlEmbedded := \"Text with \u202bembedded\u202c content\"\n\tfmt.Println(rtlEmbedded)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains PDF (Pop Directional Formatting) U+202C\n\tformatted := \"Text with \u202cformatting\"\n\tfmt.Println(formatted)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains LRO (Left-to-Right Override) U+202D\n\toverride := \"Text \u202doverride\"\n\tfmt.Println(override)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains RLO (Right-to-Left Override) U+202E\n\trloText := \"Text \u202eoverride\"\n\tfmt.Println(rloText)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains RLM (Right-to-Left Mark) U+200F\n\tmarked := \"Text \u200fmarked\"\n\tfmt.Println(marked)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Contains LRM (Left-to-Right Mark) U+200E\n\tlrmText := \"Text \u200emarked\"\n\tfmt.Println(lrmText)\n}\n"}, 1, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

// Safe code without bidirectional characters
func main() {
	username := "admin"
	password := "secret"
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

// Normal comment with regular text
func main() {
	// This is a safe comment
	isAdmin := true
	if isAdmin {
		fmt.Println("Access granted")
	}
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func main() {
	// Regular ASCII characters only
	message := "Hello, World!"
	fmt.Println(message)
}
`}, 0, gosec.NewConfig()},
		{[]string{`
package main

import "fmt"

func authenticateUser(username, password string) bool {
	// Normal authentication logic
	if username == "admin" && password == "secret" {
		return true
	}
	return false
}

func main() {
	result := authenticateUser("user", "pass")
	fmt.Println("Authenticated:", result)
}
`}, 0, gosec.NewConfig()},
	}
)
