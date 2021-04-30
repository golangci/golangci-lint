package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"

	"github.com/anduril/golangci-lint/pkg/golinters/goanalysis"
)

const (
	runtimeConfName = "runtimeconfig"
	getFn           = "Get"
	subscribeFn     = "Subscribe"
)

var (
	jsonProtoReg = regexp.MustCompile(`json=(?P<proto>[^"]*)`)
	jsonReg      = regexp.MustCompile(`json:"(?P<reg>[^"]*)`)
)

func NewRuntimeFormat() *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name:     runtimeConfName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
	return goanalysis.NewLinter(
		runtimeConfName,
		"Tool for detecting non or mis-configured runtimeConfig serialization or misuse of the API",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	i.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		target := typeutil.Callee(pass.TypesInfo, call)
		numArgs := 3
		if target == nil {
			return
		} else if target.Pkg() == nil || !(target.Pkg().Path() == "ghe.anduril.dev/anduril/graphene-go/pkg/graphene" &&
			(target.Name() == getFn || target.Name() == subscribeFn)) || len(call.Args) != numArgs {
			return
		}
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
			fnName := selector.Sel.Name
			isGetConfig := fnName == getFn
			isSubscribeConfig := fnName == subscribeFn
			if isGetConfig || isSubscribeConfig {
				defaultConfArg := call.Args[1]
				defaultConfArgType := pass.TypesInfo.Types[defaultConfArg].Type
				defaultConfig := getStructType(defaultConfArgType)
				if defaultConfig == nil {
					pass.Reportf(
						defaultConfArg.Pos(),
						"expected 2nd arg to be non nil but got %v",
						defaultConfig,
					)
					return
				}
				validateFieldSerialization(pass, defaultConfArg.Pos(), defaultConfig, defaultConfArgType.String())
				arg3 := call.Args[2]
				if isSubscribeConfig {
					validateTypeConversion(arg3, defaultConfArgType.String(), pass)
				} else if isGetConfig {
					checkPointerArg(arg3, "", defaultConfArgType.String(), arg3.Pos(), pass, false)
				}
			}
		}
	})
	return nil, nil
}

func checkPointerArg(arg ast.Expr, argName, defType string, pos token.Pos, pass *analysis.Pass, isRef bool) {
	switch t := arg.(type) {
	case *ast.Ident:
		if t.Obj == nil {
			pass.Reportf(t.Pos(), "expected configuration object as 3rd arg but got nil")
			return
		}
		if varDec, ok := t.Obj.Decl.(*ast.AssignStmt); ok {
			if len(varDec.Rhs) == 1 {
				checkPointerArg(varDec.Rhs[0], t.Name, defType, pos, pass, isRef)
			}
		}
	case *ast.UnaryExpr:
		if !t.Op.IsOperator() || t.Op.String() != "&" {
			pass.Reportf(pos, "expected ref as 3rd arg to Get")
		}
		checkPointerArg(t.X, argName, defType, pos, pass, true)
	case *ast.CompositeLit:
		if !isRef {
			pass.Reportf(pos, "expected ref as 3rd arg to Get")
		}
		if cl, ok := t.Type.(*ast.Ident); ok {
			checkSameType(cl, argName, defType, pass)
		}
	default:
		pass.Reportf(pos, "expected ref as 3rd arg to Get")
	}
}

func validateTypeConversion(arg ast.Expr, defType string, pass *analysis.Pass) {
	switch cb := arg.(type) {
	case *ast.FuncLit:
		argName := cb.Type.Params.List[0].Names[0].Name
		// opt out if the config isn't used
		if argName == "_" {
			return
		}
		for _, stmt := range cb.Body.List {
			validateTypeConversionStmt(stmt, argName, defType, pass)
		}
	case *ast.Ident:
		if cb.Obj == nil {
			pass.Reportf(cb.Pos(), "expected function as 3rd arg but got nil")
			return
		}
		if varDec, ok := cb.Obj.Decl.(*ast.AssignStmt); ok {
			if len(varDec.Rhs) == 1 {
				validateTypeConversion(varDec.Rhs[0], defType, pass)
			}
		}
	}
}

func validateTypeConversionStmt(stmt ast.Stmt, argName, defType string, pass *analysis.Pass) {
	switch t := stmt.(type) {
	case *ast.IfStmt:
		validateTypeConversionStmt(t.Init, argName, defType, pass)
		validateTypeConversionStmt(t.Body, argName, defType, pass)
		validateTypeConversionStmt(t.Else, argName, defType, pass)
	case *ast.BlockStmt:
		for _, stmt := range t.List {
			validateTypeConversionStmt(stmt, argName, defType, pass)
		}
	case *ast.AssignStmt:
		validateAssignConversion(t.Rhs[0], argName, defType, pass)
	case *ast.ExprStmt:
		if cb, ok := t.X.(*ast.CallExpr); ok {
			for _, arg := range cb.Args {
				validateAssignConversion(arg, argName, defType, pass)
			}
		}
	}
}

func validateAssignConversion(arg ast.Expr, argName, defType string, pass *analysis.Pass) {
	if expr, ok := arg.(*ast.TypeAssertExpr); ok {
		if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == argName {
			if se, ok := expr.Type.(*ast.StarExpr); !ok {
				pass.Reportf(expr.Type.Pos(), "the configuration object (arg %s) is a reference, add '*' to fix type conversion", argName)
			} else if typeConv, ok := se.X.(*ast.Ident); ok {
				checkSameType(typeConv, argName, defType, pass)
			}
		}
	}
}

func checkSameType(ident *ast.Ident, argName, defType string, pass *analysis.Pass) {
	if !strings.HasSuffix(defType, ident.Obj.Name) {
		pass.Reportf(ident.Pos(),
			"the configuration object (arg %s) has to match the type of the default argument, expected type: %s but found %s",
			argName, defType, ident.Obj.Name,
		)
	}
}

func validateFieldSerialization(pass *analysis.Pass, pos token.Pos, config *types.Struct, configTypeName string) {
	for i := 0; i < config.NumFields(); i++ {
		field := config.Field(i)
		tag := config.Tag(i)
		// oneOf's don't specify a json tag but the serializer 'correctly' uses the camelCase name
		isOneOf := strings.Contains(tag, "protobuf_oneof")

		// ignore missing tags for unexported fields and oneOfs
		if !field.Exported() || isOneOf {
			continue
		}
		ft := getStructType(field.Type())
		// recurse into complex types
		if field.Embedded() {
			validateFieldSerialization(pass, pos, ft, fmt.Sprintf("%s embedded in %s", field.Type(), configTypeName))
			continue
		}
		hasJSONTag := strings.Contains(tag, "json")

		if !hasJSONTag || !hasValidJSONTag(tag) {
			pass.Reportf(pos,
				"runtimeConfigurations must specify a json field name in camelCase format (e.g. json:'fieldName') "+
					"for exported fields, found field '%s' of '%s' using '%s'",
				field.Name(), configTypeName, tag,
			)
		}
		if ft != nil {
			validateFieldSerialization(pass, pos, ft, fmt.Sprintf("%s in %s", field.Type(), configTypeName))
		}
	}
}

func getStructType(t types.Type) *types.Struct {
	switch t := t.(type) {
	case *types.Pointer:
		return getStructType(t.Elem())
	case *types.Named:
		return getStructType(t.Underlying())
	case *types.Struct:
		return t
	default:
		return nil
	}
}

func hasValidJSONTag(tag string) bool {
	// protos contain proto and reg. json tag, proto is preferred
	matches := jsonProtoReg.FindStringSubmatch(tag)
	expectedMatches := 2
	if len(matches) != expectedMatches {
		matches = jsonReg.FindStringSubmatch(tag)
		if len(matches) != expectedMatches {
			return false
		}
	}
	for _, m := range strings.Split(matches[1], ",") {
		if m == "" || m == "omitempty" || m == "-" || m == "string" {
			continue
		}
		return camelCase(m) == m
	}
	return false
}

func camelCase(fieldName string) string {
	r := ""
	capitalize := false
	for i, c := range fieldName {
		ch := string(c)
		if capitalize {
			capitalize = false
			ch = strings.ToUpper(ch)
		}
		if i == 0 {
			ch = strings.ToLower(ch)
		} else if ch == "_" || ch == "-" {
			capitalize = true
			continue
		}
		r += ch
	}
	return r
}
