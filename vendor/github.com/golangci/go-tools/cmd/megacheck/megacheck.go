// megacheck runs staticcheck, gosimple and unused.
package megacheck // import "github.com/golangci/go-tools/cmd/megacheck"

import (
	"github.com/golangci/go-tools/lint"
	"github.com/golangci/go-tools/lint/lintutil"
	"github.com/golangci/go-tools/simple"
	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/staticcheck"
	"github.com/golangci/go-tools/unused"
	"golang.org/x/tools/go/loader"
)

func Run(program *loader.Program, conf *loader.Config, ssaProg *ssa.Program, enableStaticcheck, enableGosimple, enableUnused bool) []lint.Problem {
	var flags struct {
		staticcheck struct {
			enabled     bool
			generated   bool
			exitNonZero bool
		}
		gosimple struct {
			enabled     bool
			generated   bool
			exitNonZero bool
		}
		unused struct {
			enabled      bool
			constants    bool
			fields       bool
			functions    bool
			types        bool
			variables    bool
			debug        string
			wholeProgram bool
			reflection   bool
			exitNonZero  bool
		}
	}
	fs := lintutil.FlagSet("megacheck")
	fs.BoolVar(&flags.gosimple.enabled,
		"simple.enabled", true, "Run gosimple")
	fs.BoolVar(&flags.gosimple.generated,
		"simple.generated", false, "Check generated code")
	fs.BoolVar(&flags.gosimple.exitNonZero,
		"simple.exit-non-zero", false, "Exit non-zero if any problems were found")

	fs.BoolVar(&flags.staticcheck.enabled,
		"staticcheck.enabled", true, "Run staticcheck")
	fs.BoolVar(&flags.staticcheck.generated,
		"staticcheck.generated", false, "Check generated code (only applies to a subset of checks)")
	fs.BoolVar(&flags.staticcheck.exitNonZero,
		"staticcheck.exit-non-zero", true, "Exit non-zero if any problems were found")

	fs.BoolVar(&flags.unused.enabled,
		"unused.enabled", true, "Run unused")
	fs.BoolVar(&flags.unused.constants,
		"unused.consts", true, "Report unused constants")
	fs.BoolVar(&flags.unused.fields,
		"unused.fields", true, "Report unused fields")
	fs.BoolVar(&flags.unused.functions,
		"unused.funcs", true, "Report unused functions and methods")
	fs.BoolVar(&flags.unused.types,
		"unused.types", true, "Report unused types")
	fs.BoolVar(&flags.unused.variables,
		"unused.vars", true, "Report unused variables")
	fs.BoolVar(&flags.unused.wholeProgram,
		"unused.exported", false, "Treat arguments as a program and report unused exported identifiers")
	fs.BoolVar(&flags.unused.reflection,
		"unused.reflect", true, "Consider identifiers as used when it's likely they'll be accessed via reflection")
	fs.BoolVar(&flags.unused.exitNonZero,
		"unused.exit-non-zero", true, "Exit non-zero if any problems were found")

	flags.gosimple.enabled = enableGosimple
	flags.staticcheck.enabled = enableStaticcheck
	flags.unused.enabled = enableUnused

	var checkers []lintutil.CheckerConfig

	if flags.staticcheck.enabled {
		sac := staticcheck.NewChecker()
		sac.CheckGenerated = flags.staticcheck.generated
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker:     sac,
			ExitNonZero: flags.staticcheck.exitNonZero,
		})
	}

	if flags.gosimple.enabled {
		sc := simple.NewChecker()
		sc.CheckGenerated = flags.gosimple.generated
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker:     sc,
			ExitNonZero: flags.gosimple.exitNonZero,
		})
	}

	if flags.unused.enabled {
		var mode unused.CheckMode
		if flags.unused.constants {
			mode |= unused.CheckConstants
		}
		if flags.unused.fields {
			mode |= unused.CheckFields
		}
		if flags.unused.functions {
			mode |= unused.CheckFunctions
		}
		if flags.unused.types {
			mode |= unused.CheckTypes
		}
		if flags.unused.variables {
			mode |= unused.CheckVariables
		}
		uc := unused.NewChecker(mode)
		uc.WholeProgram = flags.unused.wholeProgram
		uc.ConsiderReflection = flags.unused.reflection
		checkers = append(checkers, lintutil.CheckerConfig{
			Checker:     unused.NewLintChecker(uc),
			ExitNonZero: flags.unused.exitNonZero,
		})

	}

	return lintutil.ProcessFlagSet(checkers, fs, program, conf, ssaProg)
}
