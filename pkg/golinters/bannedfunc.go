package golinters

import (
	"go/ast"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"gopkg.in/yaml.v2"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

// Analyzer lint 插件的结构体
var Analyzer = &analysis.Analyzer{
	Name:     "bandfunc",
	Doc:      "Checks that cannot use func",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type configSetting struct {
	LinterSettings BandFunc `yaml:"linters-settings"`
}

// BandFunc 读取配置的结构体
type BandFunc struct {
	Funcs map[string]string `yaml:"bannedfunc,flow"`
}

// NewCheckBannedFunc 返回检查函数
func NewCheckBannedFunc() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"bannedfunc",
		"Checks that cannot use func",
		[]*analysis.Analyzer{Analyzer},
		nil,
	).WithContextSetter(linterCtx).WithLoadMode(goanalysis.LoadModeSyntax)
}

func linterCtx(lintCtx *linter.Context) {
	// 读取配置文件
	config := loadConfigFile()

	configMap := configToConfigMap(config)

	Analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		useMap := getUsedMap(pass, configMap)
		for _, f := range pass.Files {
			ast.Inspect(f, astFunc(pass, useMap))
		}
		return nil, nil
	}
}

func astFunc(pass *analysis.Pass, usedMap map[string]map[string]string) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}

		m := usedMap[ident.Name]
		if m == nil {
			return true
		}

		sel := selector.Sel
		value, ok := m[sel.Name]
		if !ok {
			return true
		}
		pass.Reportf(node.Pos(), value)
		return true
	}
}

// configToConfigMap 将配置文件转成 map
// map[包名]map[函数名]错误提示
// example:
// {
//   time: {
//     Now: 不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()
//     Date: xxxx
//   },
//   github.com/MiaoSiLa/missevan-go/util/time: {
//     TimeNow: xxxxxx
//     SetTimeNow: xxxxx
//   }
// }
func configToConfigMap(config configSetting) map[string]map[string]string {
	configMap := make(map[string]map[string]string)
	for k, v := range config.LinterSettings.Funcs {
		strs := strings.Split(k, ").")
		if len(strs) != 2 {
			continue
		}
		if len(strs[0]) <= 1 || strs[0][0] != '(' {
			continue
		}
		pkg, name := strs[0][1:], strs[1]
		if name == "" {
			continue
		}
		m := configMap[pkg]
		if m == nil {
			m = make(map[string]string)
			configMap[pkg] = m
		}
		m[name] = v
	}
	return configMap
}

func loadConfigFile() configSetting {
	wd, _ := os.Getwd()
	f, err := ioutil.ReadFile(wd + "/.golangci.yml")
	if err != nil {
		panic(err)
	}
	return decodeFile(f)
}

func decodeFile(b []byte) configSetting {
	var config configSetting
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// getUsedMap 将配置文件的 map 转成文件下实际变量名的 map
// map[包的别名]map[函数名]错误提示
// example:
// {
//   time: {
//     Now: 不能使用 time.Now() 请使用 MiaoSiLa/missevan-go/util 下 TimeNow()
//     Date: xxxx
//   },
//   util: {
//     TimeNow: xxxxxx
//     SetTimeNow: xxxxx
//   }
// }
func getUsedMap(pass *analysis.Pass, configMap map[string]map[string]string) map[string]map[string]string {
	useMap := make(map[string]map[string]string)
	for _, item := range pass.Pkg.Imports() {
		if m, ok := configMap[item.Path()]; ok {
			useMap[item.Name()] = m
		}
	}
	return useMap
}
