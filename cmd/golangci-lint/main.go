package main

import (
	"log"

	"github.com/golangci/golangci-lint/internal/commands"
)

func main() {
	log.SetFlags(0) // don't print time

	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
