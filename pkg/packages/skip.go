package packages

import (
	"fmt"
	"path/filepath"
)

func pathElemRe(e string) string {
	return fmt.Sprintf(`(^|%c)%s($|%c)`, filepath.Separator, e, filepath.Separator)
}

var StdExcludeDirRegexps = []string{
	pathElemRe("vendor"),
	pathElemRe("third_party"),
	pathElemRe("testdata"),
	pathElemRe("examples"),
	pathElemRe("Godeps"),
	pathElemRe("builtin"),
}
