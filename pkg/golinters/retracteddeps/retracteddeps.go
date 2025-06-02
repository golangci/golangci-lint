// Package retracteddeps defines an analyzer that checks for retracted module
// versions in dependencies.
package retracteddeps

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.RetractedDepsSettings) *goanalysis.Linter {
	return goanalysis.NewLinterFromAnalyzer(Analyzer)
}

// Analyzer checks for retracted module versions in dependencies.
var Analyzer = &analysis.Analyzer{
	Name: "retracteddeps",
	Doc:  "check for retracted module versions in dependencies",
	Run:  run,
}

// moduleInfo holds information about a module and its retractions
type moduleInfo struct {
	Path       string       `json:"Path"`
	Version    string       `json:"Version"`
	Replace    *moduleInfo  `json:"Replace,omitempty"`
	Retracted  []string     `json:"Retracted,omitempty"`
	Deprecated string       `json:"Deprecated,omitempty"`
	Error      *moduleError `json:"Error,omitempty"`
}

type moduleError struct {
	Err string `json:"Err"`
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Find go.mod file in the module root
	goModPath, err := findGoMod(pass)
	if err != nil {
		// If no go.mod found, skip this check
		return nil, nil
	}

	// Parse the main module's go.mod
	mainMod, err := parseGoMod(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// Check each direct dependency
	for _, req := range mainMod.Require {
		if err := checkModuleRetraction(pass, req.Mod.Path, req.Mod.Version); err != nil {
			// Log error but continue checking other dependencies
			continue
		}
	}

	return nil, nil
}

// findGoMod finds the go.mod file in the module root
func findGoMod(pass *analysis.Pass) (string, error) {
	// Start from the directory of the first package
	if len(pass.Files) == 0 {
		return "", fmt.Errorf("no files to analyze")
	}

	// Get the directory of the first file
	firstFile := pass.Fset.File(pass.Files[0].Pos())
	if firstFile == nil {
		return "", fmt.Errorf("cannot determine file location")
	}

	dir := filepath.Dir(firstFile.Name())

	// Walk up the directory tree looking for go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

// parseGoMod parses a go.mod file
func parseGoMod(path string) (*modfile.File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.Parse(path, data, nil)
}

// checkModuleRetraction checks if a module version is retracted
func checkModuleRetraction(pass *analysis.Pass, modPath, modVersion string) error {
	// Try to find the module's go.mod in the cache
	modRoot := findModuleInCache(modPath, modVersion)
	if modRoot == "" {
		// Module not in cache, we can't check retractions
		return nil
	}

	goModPath := filepath.Join(modRoot, "go.mod")
	modFile, err := parseGoMod(goModPath)
	if err != nil {
		// Can't parse go.mod, skip
		return nil
	}

	// Check if there are any retractions in a newer version
	// We need to check the latest version's go.mod for retraction info
	latestModRoot := findLatestModuleInCache(modPath)
	if latestModRoot != "" && latestModRoot != modRoot {
		latestGoModPath := filepath.Join(latestModRoot, "go.mod")
		if latestModFile, err := parseGoMod(latestGoModPath); err == nil {
			// Check retractions from the latest version
			checkRetractions(pass, modPath, modVersion, latestModFile.Retract)
		}
	}

	// Also check retractions in the current version's go.mod
	checkRetractions(pass, modPath, modVersion, modFile.Retract)

	return nil
}

// checkRetractions checks if a version is covered by any retraction
func checkRetractions(pass *analysis.Pass, modPath, modVersion string, retractions []*modfile.Retract) {
	for _, retract := range retractions {
		if versionInInterval(modVersion, retract.VersionInterval) {
			msg := fmt.Sprintf("module %s@%s is retracted", modPath, modVersion)
			if retract.Rationale != "" {
				msg += ": " + retract.Rationale
			}
			pass.Reportf(reportPos(pass), msg)
			break
		}
	}
}

// findModuleInCache finds a module version in the module cache
func findModuleInCache(modPath, version string) string {
	cacheDir := getModCacheDir()
	if cacheDir == "" {
		return ""
	}

	// Escape the module path for filesystem
	escapedPath, err := module.EscapePath(modPath)
	if err != nil {
		return ""
	}

	// Clean version for directory name
	dirVersion := version
	if !strings.HasPrefix(dirVersion, "v") {
		dirVersion = "v" + dirVersion
	}

	modRoot := filepath.Join(cacheDir, escapedPath+"@"+dirVersion)
	if _, err := os.Stat(modRoot); err == nil {
		return modRoot
	}

	return ""
}

// findLatestModuleInCache tries to find the latest version of a module in cache
func findLatestModuleInCache(modPath string) string {
	cacheDir := getModCacheDir()
	if cacheDir == "" {
		return ""
	}

	escapedPath, err := module.EscapePath(modPath)
	if err != nil {
		return ""
	}
	pattern := filepath.Join(cacheDir, escapedPath+"@v*")

	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	// Find the highest version
	var latestPath string
	var latestVersion string

	for _, match := range matches {
		// Extract version from path
		base := filepath.Base(match)
		if idx := strings.LastIndex(base, "@"); idx >= 0 {
			ver := base[idx+1:]
			if latestVersion == "" || semver.Compare(ver, latestVersion) > 0 {
				latestVersion = ver
				latestPath = match
			}
		}
	}

	return latestPath
}

// getModCacheDir returns the module cache directory
func getModCacheDir() string {
	// Get GOMODCACHE or use default
	if cacheDir := os.Getenv("GOMODCACHE"); cacheDir != "" {
		return cacheDir
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// Use default GOPATH
		if home, err := os.UserHomeDir(); err == nil {
			gopath = filepath.Join(home, "go")
		}
	}

	if gopath != "" {
		return filepath.Join(gopath, "pkg", "mod")
	}

	return ""
}

// versionInInterval checks if a version is within a retraction interval
func versionInInterval(version string, interval modfile.VersionInterval) bool {
	// Clean version strings for comparison
	version = cleanVersion(version)
	low := cleanVersion(interval.Low)
	high := cleanVersion(interval.High)

	// Use semver.Compare for version comparison
	// Returns -1, 0, or 1 if version is less than, equal to, or greater than other version
	return semver.Compare(version, low) >= 0 && semver.Compare(version, high) <= 0
}

// cleanVersion ensures version strings are in canonical form
func cleanVersion(v string) string {
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return semver.Canonical(v)
}

// reportPos returns a token.Pos for reporting issues
func reportPos(pass *analysis.Pass) token.Pos {
	if len(pass.Files) > 0 {
		// Report at the package declaration
		return pass.Files[0].Name.Pos()
	}
	return token.NoPos
}
