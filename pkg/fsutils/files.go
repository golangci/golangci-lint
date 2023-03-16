package fsutils

import "path/filepath"

// Files combines different operations related to handling file paths and content.
type Files struct {
	*LineCache
	pathPrefix string
}

func NewFiles(lc *LineCache, pathPrefix string) *Files {
	return &Files{
		LineCache:  lc,
		pathPrefix: pathPrefix,
	}
}

// WithPathPrefix takes a path that is relative to the current directory (as used in issues)
// and adds the configured path prefix, if there is one.
// The resulting path then can be shown to the user or compared against paths specified in the configuration.
func (f *Files) WithPathPrefix(relativePath string) string {
	return WithPathPrefix(f.pathPrefix, relativePath)
}

// WithPathPrefix takes a path that is relative to the current directory (as used in issues)
// and adds the configured path prefix, if there is one.
// The resulting path then can be shown to the user or compared against paths specified in the configuration.
func WithPathPrefix(pathPrefix, relativePath string) string {
	if pathPrefix == "" {
		return relativePath
	}
	return filepath.Join(pathPrefix, relativePath)
}
