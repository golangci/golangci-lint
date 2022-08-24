//go:build !windows

package processors

// normalizePathInRegex it's a noop function on Unix.
func normalizePathInRegex(path string) string {
	return path
}
