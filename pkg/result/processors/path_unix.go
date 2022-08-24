//go:build !windows

package processors

func normalizePathInRegex(path string) string {
	return path
}
