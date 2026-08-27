package cache

import (
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// RelativePosition creates a [token.Position] with a filename relative to the module directory.
func RelativePosition(mod *packages.Module, pos token.Position) token.Position {
	pos.Filename, _ = toRelativePath(mod, pos.Filename)

	return pos
}

// AbsolutionPosition creates the absolute path base on the module directory.
func AbsolutionPosition(mod *packages.Module, pos token.Position) token.Position {
	pos.Filename = toAbsolutePath(mod, pos.Filename)

	return pos
}

func toRelativePath(mod *packages.Module, filename string) (string, bool) {
	if filename == "" || !isCurrentModule(mod) {
		return filename, false
	}

	rel, err := filepath.Rel(mod.Dir, filename)
	if err != nil {
		return filename, false
	}

	rel = filepath.ToSlash(rel)

	if rel == ".." || strings.HasPrefix(rel, "../") {
		return filename, false
	}

	return rel, true
}

func toAbsolutePath(mod *packages.Module, filename string) string {
	if filename == "" || !isCurrentModule(mod) || filepath.IsAbs(filename) {
		return filename
	}

	return filepath.Join(mod.Dir, filepath.FromSlash(filename))
}

func isCurrentModule(mod *packages.Module) bool {
	return mod != nil && mod.Dir != "" && mod.Version == ""
}
