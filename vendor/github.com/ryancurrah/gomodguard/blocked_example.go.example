package gomodguard

import (
	"io/ioutil"

	"github.com/gofrs/uuid"
	module "github.com/uudashr/go-module"
)

func aBlockedImport() { // nolint: deadcode,unused
	b, err := ioutil.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	mod, err := module.Parse(b)
	if err != nil {
		panic(err)
	}

	_ = mod

	_ = uuid.Must(uuid.NewV4())
}
