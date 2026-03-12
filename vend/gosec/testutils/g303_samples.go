package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG303 - bad tempfile permissions & hardcoded shared path
var SampleCodeG303 = []CodeSample{
	{[]string{`
package samples

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func main() {
	err := ioutil.WriteFile("/tmp/demo2", []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	f, err := os.Create("/tmp/demo2")
	if err != nil {
		fmt.Println("Error while writing!")
	} else if err = f.Close(); err != nil {
		fmt.Println("Error while closing!")
	}
	err = os.WriteFile("/tmp/demo2", []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile("/usr/tmp/demo2", []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile("/tmp/" + "demo2", []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile(os.TempDir() + "/demo2", []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile(path.Join("/var/tmp", "demo2"), []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile(path.Join(os.TempDir(), "demo2"), []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
	err = os.WriteFile(filepath.Join(os.TempDir(), "demo2"), []byte("This is some data"), 0644)
	if err != nil {
		fmt.Println("Error while writing!")
	}
}
`}, 9, gosec.NewConfig()},
}
