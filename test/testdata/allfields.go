//golangcitest:args -Eallfields
package testdata

import "time"

type user struct {
	Name      string
	CreatedAt time.Time
}

func Allfields() {
	_ = user{
		Name:      "John",
		CreatedAt: time.Now(),
		//allfields
	}
	_ = user{ // want "field CreatedAt is not set"
		Name: "John",
		//allfields
	}
	_ = user{
		Name: "John",
	}
}
