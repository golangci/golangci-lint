package golinters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNewBannedFunc(t *testing.T) {
	lint := NewBannedFunc(&config.BannedFuncSettings{Funcs: map[string]string{"(time).Now": "Disable time.Now"}})
	if lint == nil {
		t.Fatal("expected lint to be not nil")
	}
}

type errorfunc func(string)

func (f errorfunc) Errorf(format string, args ...interface{}) {
	f(fmt.Sprintf(format, args...))
}

func TestBannedFunc_Run(t *testing.T) {
	files := map[string]string{
		"a/b.go": `
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello")
	_ = time.Now()
	_ = time.Now().Unix()
	_ = time.Unix(10000,0)
}
`,
	}
	dir, cleanup, err := analysistest.WriteFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	var got []string
	fakeT := errorfunc(func(s string) { got = append(got, s) }) // a fake *testing.T
	bf := &bannedFunc{
		ban: &config.BannedFuncSettings{
			Funcs: map[string]string{
				"(time).Now":    "Disable time.Now",
				"(time).Unix":   "Disable time.Unix",
				"(fmt).Println": "Disable fmt.Println",
			},
		},
	}
	analysistest.Run(fakeT, dir, &analysis.Analyzer{Run: bf.Run}, "a")
	want := []string{
		`a/b.go:9:14: unexpected diagnostic: Disable fmt.Println`,
		`a/b.go:10:6: unexpected diagnostic: Disable time.Now`,
		`a/b.go:11:6: unexpected diagnostic: Disable time.Now`,
		`a/b.go:12:6: unexpected diagnostic: Disable time.Unix`,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBannedFunc_parseBannedFunc(t *testing.T) {
	bf := &bannedFunc{
		ban: &config.BannedFuncSettings{
			Funcs: map[string]string{
				"(ioutil).WriteFile":              "As of Go 1.16, this function simply calls os.WriteFile.",
				"(ioutil).ReadFile":               "As of Go 1.16, this function simply calls os.ReadFile.",
				"(github.com/example/banned).New": "This function is deprecated",
				"(github.com/example/banned).":    "Skip checking for empty function names",
				"().":                             "Empty",
				").":                              "Empty",
			},
		},
	}

	confMap := bf.parseBannedFunc()
	if len(confMap) != 2 {
		t.Fatalf("expected 2, got %d", len(confMap))
	}
	if len(confMap["ioutil"]) != 2 {
		t.Fatalf("expected 2, got %d", len(confMap["ioutil"]))
	}
	if confMap["ioutil"]["WriteFile"] != "As of Go 1.16, this function simply calls os.WriteFile." {
		t.Errorf("expected 'As of Go 1.16, this function simply calls os.WriteFile.', got %s", confMap["ioutil"]["WriteFile"])
	}
	if confMap["ioutil"]["ReadFile"] != "As of Go 1.16, this function simply calls os.ReadFile." {
		t.Errorf("expected 'As of Go 1.16, this function simply calls os.ReadFile.', got %s", confMap["ioutil"]["ReadFile"])
	}
	if len(confMap["github.com/example/banned"]) != 1 {
		t.Fatalf("expected 1, got %d", len(confMap["github.com/example/banned"]))
	}
	if confMap["github.com/example/banned"]["New"] != "This function is deprecated" {
		t.Errorf("expected 'This function is deprecated', got %s", confMap["github.com/example/banned"]["New"])
	}
}
