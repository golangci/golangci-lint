// nolint:dupl
package golinters

import (
	"go/ast"
	"go/token"
	"go/types"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
	"github.com/anduril/golangci-lint/pkg/lint/linter"
	"github.com/anduril/golangci-lint/pkg/result"
)

const testifyAssertEqualProto = "testifyeqproto"

func NewTestifyAssertEqualProto() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     testifyAssertEqualProto,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	return goanalysis.NewLinter(
		testifyAssertEqualProto,
		"Finds places in which assert.Equal is invoked on two structs that transitively reference a proto message",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues := runTestifyEqProto(pass)

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runTestifyEqProto(pass *analysis.Pass) []goanalysis.Issue {
	var issues []goanalysis.Issue
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	i.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return // not a static call
		} else if fn.FullName() != testifyAssertEqualMethod {
			return
		} else if len(call.Args) < 3 { // nolint:gomnd
			return
		}

		type2 := pass.TypesInfo.Types[call.Args[1]].Type
		type3 := pass.TypesInfo.Types[call.Args[2]].Type

		visited := map[string]bool{}
		if transitivelyReferencesProto(type2, visited) || transitivelyReferencesProto(type3, visited) {
			issue := goanalysis.NewIssue(&result.Issue{
				Pos: pass.Fset.Position(n.Pos()),
				Text: "call to assert.Equal made with structs that contain proto.Message fields: " +
					"Use proto.Equal instead.",
			}, pass)
			issues = append(issues, issue)
		}
	})
	return issues
}

func transitivelyReferencesProto(t types.Type, visited map[string]bool) bool {
	switch t := t.(type) {
	case *types.Pointer:
		return transitivelyReferencesProto(t.Elem(), visited)
	case *types.Named:
		if types.Implements(types.NewPointer(t), protoMsgType) {
			return true
		}
		if visited[t.String()] {
			return false
		}
		visited[t.String()] = true
		return transitivelyReferencesProto(t.Underlying(), visited)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			if transitivelyReferencesProto(f.Type(), visited) {
				return true
			}
		}
		return false
	case *types.Interface:
		return false
	default:
		return false
	}
}

var protoMsgType *types.Interface

// Construct a proto.Message interface type:
//
// type MessageV1 interface {
//   Reset()
//   String() string
//   ProtoMessage()
// }
func init() { // nolint:gochecknoinits
	nullary := types.NewSignature(nil, nil, nil, false) // func()
	retstr := types.NewSignature(                       // func() string
		nil,
		nil,
		types.NewTuple(types.NewParam(token.NoPos, nil, "", types.Typ[types.String])),
		false,
	)
	methods := []*types.Func{
		types.NewFunc(token.NoPos, nil, "Reset", nullary),
		types.NewFunc(token.NoPos, nil, "String", retstr),
		types.NewFunc(token.NoPos, nil, "ProtoMessage", nullary),
	}
	protoMsgType = types.NewInterfaceType(methods, nil).Complete()
}
