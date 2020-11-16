package golinters

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDecodeFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	b := strings.TrimSpace(`
linters-settings:
  bannedfunc:
    (time).Now: "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()"
    (github.com/MiaoSiLa/missevan-go/util).TimeNow: "aaa"
`)
	var setting configSetting
	require.NotPanics(func() { setting = decodeFile([]byte(b)) })
	require.NotNil(setting.LinterSettings)
	val := setting.LinterSettings.Funcs["(time).Now"]
	assert.NotEmpty(val)
	val = setting.LinterSettings.Funcs["(github.com/MiaoSiLa/missevan-go/util).TimeNow"]
	assert.NotEmpty(val)
}

func TestConfigToConfigMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := map[string]string{
		"(time).Now": "不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()",
		"(github.com/MiaoSiLa/missevan-go/util).TimeNow": "xxxx",
		"().": "(). 情况",
		").":  "). 情况",
	}
	s := configSetting{LinterSettings: BandFunc{Funcs: m}}
	setting := configToConfigMap(s)
	require.Len(setting, 2)
	require.NotNil(setting["time"])
	require.NotNil(setting["time"]["Now"])
	assert.Equal("不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()", setting["time"]["Now"])
	require.NotNil(setting["github.com/MiaoSiLa/missevan-go/util"])
	require.NotNil(setting["github.com/MiaoSiLa/missevan-go/util"]["TimeNow"])
	assert.Equal("xxxx", setting["github.com/MiaoSiLa/missevan-go/util"]["TimeNow"])
	assert.Nil(setting["()."])
	assert.Nil(setting[")."])
}

func TestGetUsedMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pkg := types.NewPackage("test", "test")
	importPkg := []*types.Package{types.NewPackage("time", "time"),
		types.NewPackage("github.com/MiaoSiLa/missevan-go/util", "util")}
	pkg.SetImports(importPkg)
	pass := analysis.Pass{Pkg: pkg}
	m := map[string]map[string]string{
		"time":                                 {"Now": "xxxx", "Date": "xxxx"},
		"assert":                               {"New": "xxxx"},
		"github.com/MiaoSiLa/missevan-go/util": {"TimeNow": "xxxx"},
	}
	usedMap := getUsedMap(&pass, m)
	require.Len(usedMap, 2)
	require.Len(usedMap["time"], 2)
	assert.NotEmpty(usedMap["time"]["Now"])
	assert.NotEmpty(usedMap["time"]["Date"])
	require.Len(usedMap["util"], 1)
	assert.NotEmpty(usedMap["util"]["TimeNow"])
}

func TestAstFunc(t *testing.T) {
	assert := assert.New(t)

	// 初始化测试用参数
	var testStr string
	pass := analysis.Pass{Report: func(diagnostic analysis.Diagnostic) {
		testStr = diagnostic.Message
	}}
	m := map[string]map[string]string{
		"time": {"Now": "time.Now"},
		"util": {"TimeNow": "util.TimeNow"},
	}
	f := astFunc(&pass, m)

	// 测试不符合情况
	node := ast.SelectorExpr{X: &ast.Ident{Name: "time"}, Sel: &ast.Ident{Name: "Date"}}
	f(&node)
	assert.Empty(testStr)
	node = ast.SelectorExpr{X: &ast.Ident{Name: "assert"}, Sel: &ast.Ident{Name: "New"}}
	f(&node)
	assert.Empty(testStr)

	// 测试符合情况
	node = ast.SelectorExpr{X: &ast.Ident{Name: "time"}, Sel: &ast.Ident{Name: "Now"}}
	f(&node)
	assert.Equal("time.Now", testStr)
	node = ast.SelectorExpr{X: &ast.Ident{Name: "util"}, Sel: &ast.Ident{Name: "TimeNow"}}
	f(&node)
	assert.Equal("util.TimeNow", testStr)
}

func TestFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	b := strings.TrimSpace(`
linters-settings:
  bannedfunc:
    (time).Now: "禁止使用 time.Now"
    (time).Unix: "禁止使用 time.Unix"
`)
	var setting configSetting
	require.NotPanics(func() { setting = decodeFile([]byte(b)) })
	var m map[string]map[string]string
	m = configToConfigMap(setting)
	Analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		useMap := getUsedMap(pass, m)
		for _, f := range pass.Files {
			ast.Inspect(f, astFunc(pass, useMap))
		}
		return nil, nil
	}
	// NOTICE: 因为配置的初始化函数将 GOPATH 设置为当前目录
	// 所以只能 import 内置包和当前目录下的包
	files := map[string]string{
		"a/b.go": `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now())
	_ = time.Now()
	_ = time.Now().Unix()
	_ = time.Unix(10000,0)
}

`}
	dir, cleanup, err := analysistest.WriteFiles(files)
	require.NoError(err)
	defer cleanup()
	var got []string
	t2 := errorfunc(func(s string) { got = append(got, s) }) // a fake *testing.T
	analysistest.Run(t2, dir, Analyzer, "a")
	want := []string{
		`a/b.go:9:14: unexpected diagnostic: 禁止使用 time.Now`,
		`a/b.go:10:6: unexpected diagnostic: 禁止使用 time.Now`,
		`a/b.go:11:6: unexpected diagnostic: 禁止使用 time.Now`,
		`a/b.go:12:6: unexpected diagnostic: 禁止使用 time.Unix`,
	}
	assert.Equal(want, got)
}

type errorfunc func(string)

func (f errorfunc) Errorf(format string, args ...interface{}) {
	f(fmt.Sprintf(format, args...))
}
