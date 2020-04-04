package gomodguard

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	blockedReasonNotInAllowedList = "import of package `%s` is blocked because the module is not in the allowed modules list."
	blockedReasonInBlockedList    = "import of package `%s` is blocked because the module is in the blocked modules list."
	goModFilename                 = "go.mod"
)

// Recommendations are alternative modules to use and a reason why.
type Recommendations struct {
	Recommendations []string `yaml:"recommendations"`
	Reason          string   `yaml:"reason"`
}

// String returns the recommended modules and reason message.
func (r *Recommendations) String() string {
	msg := ""

	if r == nil {
		return msg
	}

	for i := range r.Recommendations {
		switch {
		case len(r.Recommendations) == 1:
			msg += fmt.Sprintf("`%s` is a recommended module.", r.Recommendations[i])
		case (i+1) != len(r.Recommendations) && (i+1) == (len(r.Recommendations)-1):
			msg += fmt.Sprintf("`%s` ", r.Recommendations[i])
		case (i + 1) != len(r.Recommendations):
			msg += fmt.Sprintf("`%s`, ", r.Recommendations[i])
		default:
			msg += fmt.Sprintf("and `%s` are recommended modules.", r.Recommendations[i])
		}
	}

	if r.Reason != "" {
		msg += fmt.Sprintf(" %s", r.Reason)
	}

	return msg
}

// HasRecommendations returns true if the blocked package has
// recommended modules.
func (r *Recommendations) HasRecommendations(pkg string) bool {
	return len(r.Recommendations) > 0
}

// BlockedModule is a blocked module name and
// optionally a list of recommended modules
// and a reason message.
type BlockedModule map[string]Recommendations

// BlockedModules a list of blocked modules.
type BlockedModules []BlockedModule

// Get returns the modules that are blocked.
func (b BlockedModules) Get() []string {
	modules := make([]string, len(b))

	for i := range b {
		for module := range b[i] {
			modules[i] = module
			break
		}
	}

	return modules
}

// RecommendedModules will return a list of recommended modules for the
// package provided. If there is no recommendation nil will be returned.
func (b BlockedModules) RecommendedModules(pkg string) *Recommendations {
	for i := range b {
		for blockedModule, recommendations := range b[i] {
			if strings.HasPrefix(strings.ToLower(pkg), strings.ToLower(blockedModule)) && recommendations.HasRecommendations(pkg) {
				return &recommendations
			}

			break
		}
	}

	return nil
}

// IsBlockedPackage returns true if the package name is in
// the blocked modules list.
func (b BlockedModules) IsBlockedPackage(pkg string) bool {
	blockedModules := b.Get()
	for i := range blockedModules {
		if strings.HasPrefix(strings.ToLower(pkg), strings.ToLower(blockedModules[i])) {
			return true
		}
	}

	return false
}

// IsBlockedModule returns true if the given module name is in the
// blocked modules list.
func (b BlockedModules) IsBlockedModule(module string) bool {
	blockedModules := b.Get()
	for i := range blockedModules {
		if strings.EqualFold(module, strings.TrimSpace(blockedModules[i])) {
			return true
		}
	}

	return false
}

// Allowed is a list of modules and module
// domains that are allowed to be used.
type Allowed struct {
	Modules []string `yaml:"modules"`
	Domains []string `yaml:"domains"`
}

// IsAllowedModule returns true if the given module
// name is in the allowed modules list.
func (a *Allowed) IsAllowedModule(module string) bool {
	allowedModules := a.Modules
	for i := range allowedModules {
		if strings.EqualFold(module, strings.TrimSpace(allowedModules[i])) {
			return true
		}
	}

	return false
}

// IsAllowedModuleDomain returns true if the given modules domain is
// in the allowed module domains list.
func (a *Allowed) IsAllowedModuleDomain(module string) bool {
	allowedDomains := a.Domains
	for i := range allowedDomains {
		if strings.HasPrefix(strings.ToLower(module), strings.TrimSpace(strings.ToLower(allowedDomains[i]))) {
			return true
		}
	}

	return false
}

// Blocked is a list of modules that are
// blocked and not to be used.
type Blocked struct {
	Modules BlockedModules `yaml:"modules"`
}

// Configuration of gomodguard allow and block lists.
type Configuration struct {
	Allowed Allowed `yaml:"allowed"`
	Blocked Blocked `yaml:"blocked"`
}

// Result represents the result of one error.
type Result struct {
	FileName   string
	LineNumber int
	Position   token.Position
	Reason     string
}

// String returns the filename, line
// number and reason of a Result.
func (r *Result) String() string {
	return fmt.Sprintf("%s:%d: %s", r.FileName, r.LineNumber, r.Reason)
}

// Processor processes Go files.
type Processor struct {
	config                    Configuration
	logger                    *log.Logger
	modfile                   *modfile.File
	blockedModulesFromModFile []string
	result                    []Result
}

// NewProcessor will create a Processor to lint blocked packages.
func NewProcessor(config Configuration, logger *log.Logger) (*Processor, error) {
	goModFileBytes, err := loadGoModFile()
	if err != nil {
		errMsg := fmt.Sprintf("unable to read %s file: %s", goModFilename, err)

		return nil, fmt.Errorf(errMsg)
	}

	mfile, err := modfile.Parse(goModFilename, goModFileBytes, nil)
	if err != nil {
		errMsg := fmt.Sprintf("unable to parse %s file: %s", goModFilename, err)

		return nil, fmt.Errorf(errMsg)
	}

	logger.Printf("info: allowed modules, %+v", config.Allowed.Modules)
	logger.Printf("info: allowed module domains, %+v", config.Allowed.Domains)
	logger.Printf("info: blocked modules, %+v", config.Blocked.Modules.Get())

	p := &Processor{
		config:  config,
		logger:  logger,
		modfile: mfile,
		result:  []Result{},
	}

	p.setBlockedModulesFromModFile()

	return p, nil
}

// ProcessFiles takes a string slice with file names (full paths)
// and lints them.
func (p *Processor) ProcessFiles(filenames []string) []Result {
	pluralModuleMsg := "s"
	if len(p.blockedModulesFromModFile) == 1 {
		pluralModuleMsg = ""
	}

	p.logger.Printf("info: found `%d` blocked module%s in the %s file, %+v",
		len(p.blockedModulesFromModFile), pluralModuleMsg, goModFilename, p.blockedModulesFromModFile)

	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			p.result = append(p.result, Result{
				FileName:   filename,
				LineNumber: 0,
				Reason:     fmt.Sprintf("unable to read file, file cannot be linted (%s)", err.Error()),
			})
		}

		p.process(filename, data)
	}

	return p.result
}

// process file imports and add lint error if blocked package is imported.
func (p *Processor) process(filename string, data []byte) {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, filename, data, parser.ParseComments)
	if err != nil {
		p.result = append(p.result, Result{
			FileName:   filename,
			LineNumber: 0,
			Reason:     fmt.Sprintf("invalid syntax, file cannot be linted (%s)", err.Error()),
		})

		return
	}

	imports := file.Imports
	for i := range imports {
		importedPkg := strings.TrimSpace(strings.Trim(imports[i].Path.Value, "\""))
		if p.isBlockedPackageFromModFile(importedPkg) {
			reason := ""

			if p.config.Blocked.Modules.IsBlockedPackage(importedPkg) {
				reason = fmt.Sprintf(blockedReasonInBlockedList, importedPkg)
			} else {
				reason = fmt.Sprintf(blockedReasonNotInAllowedList, importedPkg)
			}

			recommendedModules := p.config.Blocked.Modules.RecommendedModules(importedPkg)
			if recommendedModules != nil {
				reason += fmt.Sprintf(" %s", recommendedModules.String())
			}

			p.addError(fileSet, imports[i].Pos(), reason)
		}
	}
}

// addError adds an error for the file and line number for the current token.Pos
// with the given reason.
func (p *Processor) addError(fileset *token.FileSet, pos token.Pos, reason string) {
	position := fileset.Position(pos)

	p.result = append(p.result, Result{
		FileName:   position.Filename,
		LineNumber: position.Line,
		Position:   position,
		Reason:     reason,
	})
}

// setBlockedModules determines which modules are blocked by reading
// the go.mod file and comparing the require modules to the allowed modules.
func (p *Processor) setBlockedModulesFromModFile() {
	blockedModules := make([]string, 0, len(p.modfile.Require))
	requiredModules := p.modfile.Require

	for i := range requiredModules {
		if !requiredModules[i].Indirect {
			requiredModule := strings.TrimSpace(requiredModules[i].Mod.Path)

			if p.config.Allowed.IsAllowedModuleDomain(requiredModule) {
				continue
			}

			if p.config.Allowed.IsAllowedModule(requiredModule) {
				continue
			}

			if len(p.config.Allowed.Modules) == 0 &&
				len(p.config.Allowed.Domains) == 0 &&
				!p.config.Blocked.Modules.IsBlockedModule(requiredModule) {
				continue
			}

			blockedModules = append(blockedModules, requiredModule)
		}
	}

	if len(blockedModules) > 0 {
		p.blockedModulesFromModFile = blockedModules
	}
}

// isBlockedPackageFromModFile returns true if the imported packages
// module is in the go.mod file and was blocked.
func (p *Processor) isBlockedPackageFromModFile(pkg string) bool {
	blockedModulesFromModFile := p.blockedModulesFromModFile
	for i := range blockedModulesFromModFile {
		if strings.HasPrefix(strings.ToLower(pkg), strings.ToLower(blockedModulesFromModFile[i])) {
			return true
		}
	}

	return false
}

func loadGoModFile() ([]byte, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.GoMod}}")
	stdout, _ := cmd.StdoutPipe()
	_ = cmd.Start()

	goModFileLocation := ""

	if stdout != nil {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(stdout)
		goModFileLocation = strings.TrimSpace(buf.String())
	}

	if _, err := os.Stat(goModFileLocation); os.IsNotExist(err) {
		return ioutil.ReadFile(goModFilename)
	}

	return ioutil.ReadFile(goModFileLocation)
}
