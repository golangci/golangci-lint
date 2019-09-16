// Package nfpm provides ways to package programs in some linux packaging
// formats.
package nfpm

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/imdario/mergo"

	"gopkg.in/yaml.v2"
)

// nolint: gochecknoglobals
var (
	packagers = map[string]Packager{}
	lock      sync.Mutex
)

// Register a new packager for the given format
func Register(format string, p Packager) {
	lock.Lock()
	packagers[format] = p
	lock.Unlock()
}

// Get a packager for the given format
func Get(format string) (Packager, error) {
	p, ok := packagers[format]
	if !ok {
		return nil, fmt.Errorf("no packager registered for the format %s", format)
	}
	return p, nil
}

// Parse decodes YAML data from an io.Reader into a configuration struct
func Parse(in io.Reader) (config Config, err error) {
	dec := yaml.NewDecoder(in)
	dec.SetStrict(true)
	if err = dec.Decode(&config); err != nil {
		return
	}

	config.Info.Version = os.ExpandEnv(config.Info.Version)
	err = config.Validate()
	return
}

// ParseFile decodes YAML data from a file path into a configuration struct
func ParseFile(path string) (config Config, err error) {
	var file *os.File
	file, err = os.Open(path) //nolint:gosec
	if err != nil {
		return
	}
	defer file.Close() // nolint: errcheck
	return Parse(file)
}

// Packager represents any packager implementation
type Packager interface {
	Package(info Info, w io.Writer) error
}

// Config contains the top level configuration for packages
type Config struct {
	Info      `yaml:",inline"`
	Overrides map[string]Overridables `yaml:"overrides,omitempty"`
}

// Get returns the Info struct for the given packager format. Overrides
// for the given format are merged into the final struct
func (c *Config) Get(format string) (info Info, err error) {
	// make a deep copy of info
	if err = mergo.Merge(&info, c.Info); err != nil {
		return
	}
	override, ok := c.Overrides[format]
	if !ok {
		// no overrides
		return
	}
	err = mergo.Merge(&info.Overridables, override, mergo.WithOverride)
	if err != nil {
		return
	}
	return
}

// Validate ensures that the config is well typed
func (c *Config) Validate() error {
	for format := range c.Overrides {
		if _, err := Get(format); err != nil {
			return err
		}
	}
	return nil
}

// Info contains information about a single package
type Info struct {
	Overridables `yaml:",inline"`
	Name         string `yaml:"name,omitempty"`
	Arch         string `yaml:"arch,omitempty"`
	Platform     string `yaml:"platform,omitempty"`
	Epoch        string `yaml:"epoch,omitempty"`
	Version      string `yaml:"version,omitempty"`
	Section      string `yaml:"section,omitempty"`
	Priority     string `yaml:"priority,omitempty"`
	Maintainer   string `yaml:"maintainer,omitempty"`
	Description  string `yaml:"description,omitempty"`
	Vendor       string `yaml:"vendor,omitempty"`
	Homepage     string `yaml:"homepage,omitempty"`
	License      string `yaml:"license,omitempty"`
	Bindir       string `yaml:"bindir,omitempty"`
}

// Overridables contain the field which are overridable in a package
type Overridables struct {
	Replaces     []string          `yaml:"replaces,omitempty"`
	Provides     []string          `yaml:"provides,omitempty"`
	Depends      []string          `yaml:"depends,omitempty"`
	Recommends   []string          `yaml:"recommends,omitempty"`
	Suggests     []string          `yaml:"suggests,omitempty"`
	Conflicts    []string          `yaml:"conflicts,omitempty"`
	Files        map[string]string `yaml:"files,omitempty"`
	ConfigFiles  map[string]string `yaml:"config_files,omitempty"`
	EmptyFolders []string          `yaml:"empty_folders,omitempty"`
	Scripts      Scripts           `yaml:"scripts,omitempty"`
	RPM          RPM               `yaml:"rpm,omitempty"`
	Deb          Deb               `yaml:"deb,omitempty"`
}

// RPM is custom configs that are only available on RPM packages
type RPM struct {
	Group       string `yaml:"group,omitempty"`
	Prefix      string `yaml:"prefix,omitempty"`
	Compression string `yaml:"compression,omitempty"`
	Release     string `yaml:"release,omitempty"`
}

// Deb is custom configs that are only available on deb packages
type Deb struct {
	Scripts DebScripts `yaml:"scripts,omitempty"`
}

// DebScripts is scripts only available on deb packages
type DebScripts struct {
	Rules string `yaml:"rules,omitempty"`
}

// Scripts contains information about maintainer scripts for packages
type Scripts struct {
	PreInstall  string `yaml:"preinstall,omitempty"`
	PostInstall string `yaml:"postinstall,omitempty"`
	PreRemove   string `yaml:"preremove,omitempty"`
	PostRemove  string `yaml:"postremove,omitempty"`
}

// Validate the given Info and returns an error if it is invalid.
func Validate(info Info) error {
	if info.Name == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if info.Arch == "" {
		return fmt.Errorf("package arch must be provided")
	}
	if info.Version == "" {
		return fmt.Errorf("package version must be provided")
	}
	if len(info.Files)+len(info.ConfigFiles) == 0 {
		return fmt.Errorf("no files were provided")
	}
	return nil
}

// WithDefaults set some sane defaults into the given Info
func WithDefaults(info Info) Info {
	if info.Bindir == "" {
		info.Bindir = "/usr/local/bin"
	}
	if info.Platform == "" {
		info.Platform = "linux"
	}
	if info.Description == "" {
		info.Description = "no description given"
	}
	info.Version = strings.TrimPrefix(info.Version, "v")
	return info
}
