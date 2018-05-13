package main

import (
	"log"

	"github.com/golangci/golangci-lint/pkg/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	log.SetFlags(0) // don't print time
	logrus.SetLevel(logrus.WarnLevel)

	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
