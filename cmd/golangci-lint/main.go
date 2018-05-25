package main

import (
	"github.com/golangci/golangci-lint/pkg/commands"
)

func main() {
	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
