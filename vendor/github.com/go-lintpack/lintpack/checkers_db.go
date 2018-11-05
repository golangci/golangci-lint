package lintpack

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"github.com/go-toolsmith/astfmt"
)

// prototypes is a set of registered checkers that are not yet instantiated.
// Registration should be done with AddChecker function.
// Initialized checkers can be obtained with NewChecker function.
var prototypes = make(map[string]checkerProto)

// GetCheckersInfo returns a checkers info list for all registered checkers.
// The slice is sorted by a checker name.
//
// Info objects can be used to instantiate checkers with NewChecker function.
func GetCheckersInfo() []*CheckerInfo {
	infoList := make([]*CheckerInfo, 0, len(prototypes))
	for _, proto := range prototypes {
		infoCopy := *proto.info
		infoList = append(infoList, &infoCopy)
	}
	sort.Slice(infoList, func(i, j int) bool {
		return infoList[i].Name < infoList[j].Name
	})
	return infoList
}

// NewChecker returns initialized checker identified by an info.
// info must be non-nil.
// Panics if info describes a checker that was not properly registered.
//
// params argument specifies per-checker options.NewChecker. Can be nil.
func NewChecker(ctx *Context, info *CheckerInfo, params map[string]interface{}) *Checker {
	proto, ok := prototypes[info.Name]
	if !ok {
		panic(fmt.Sprintf("checker with name %q not registered", info.Name))
	}
	return proto.constructor(ctx, params)
}

// FileWalker is an interface every checker should implement.
//
// The WalkFile method is executed for every Go file inside the
// package that is being checked.
type FileWalker interface {
	WalkFile(*ast.File)
}

// AddChecker registers a new checker into a checkers pool.
// Constructor is used to create a new checker instance.
// Checker name (defined in CheckerInfo.Name) must be unique.
//
// If checker is never needed, for example if it is disabled,
// constructor will not be called.
func AddChecker(info *CheckerInfo, constructor func(*CheckerContext) FileWalker) {
	if _, ok := prototypes[info.Name]; ok {
		panic(fmt.Sprintf("checker with name %q already registered", info.Name))
	}

	trimDocumentation := func(d *CheckerInfo) {
		fields := []*string{
			&d.Summary,
			&d.Details,
			&d.Before,
			&d.After,
			&d.Note,
		}
		for _, f := range fields {
			*f = strings.TrimSpace(*f)
		}
	}
	validateDocumentation := func(d *CheckerInfo) {
		// TODO(Quasilyte): validate documentation.
	}

	trimDocumentation(info)
	validateDocumentation(info)

	proto := checkerProto{
		info: info,
		constructor: func(ctx *Context, params parameters) *Checker {
			var c Checker
			c.Info = info
			c.ctx = CheckerContext{
				Context: ctx,
				Params:  params,
				printer: astfmt.NewPrinter(ctx.FileSet),
			}
			c.fileWalker = constructor(&c.ctx)
			return &c
		},
	}

	prototypes[info.Name] = proto
}
