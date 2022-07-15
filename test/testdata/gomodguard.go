//golangcitest:args -Egomodguard
//golangcitest:config_path testdata/configs/gomodguard.yml
package testdata

import (
	"log"

	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3" // ERROR "import of package `gopkg.in/yaml.v3` is blocked because the module is in the blocked modules list. `github.com/kylelemons/go-gypsy` is a recommended module. This is an example of recommendations."
)

// Something just some struct
type Something struct{}

func aAllowedImport() { // nolint: deadcode,unused
	mfile, _ := modfile.Parse("go.mod", []byte{}, nil)

	log.Println(mfile)
}

func aBlockedImport() { // nolint: deadcode,unused
	data := []byte{}
	something := Something{}
	_ = yaml.Unmarshal(data, &something)

	log.Println(data)
}
